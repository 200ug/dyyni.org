import lottie from "lottie-web"
import pako from "pako"

const TGS_DIR_PREFIX = "/static/dice"
const DIE_ANIMATION_COUNT = 6
const ROLL_RESULTS = [
    {
        icon: "/static/icons/email.svg",
        alt: "Email",
        text: "a (at) this domain"
    },
    {
        icon: "/static/icons/pgp.svg",
        alt: "PGP key",
        link: "/static/pgp.asc",
        text: "1530 .. 2B3E"
    },
    {
        icon: "/static/icons/signal.svg",
        alt: "Signal",
        link: "https://signal.me/#eu/bsHM3c8pZCwBgZqPF4V2AmDpMww7K2bxD6iAfECyj6swnfG6v7hBbZYYJOPxiuky",
        text: "mix.99"
    },
    {
        icon: "/static/icons/monero.svg",
        alt: "Monero",
        link: "/static/monero.txt",
        text: "83B5pKo .. FHe23Y"
    },
    {
        icon: "/static/icons/telegram.svg",
        alt: "Telegram",
        link: "https://t.me/japangrab",
        text: "japangrab"
    },
    {
        icon: "/static/icons/codeberg.svg",
        alt: "Codeberg",
        link: "https://codeberg.org/2ug",
        text: "2ug"
    }
]

// null (not loaded) | object (decompressed lottie json)
const animData = new Array(DIE_ANIMATION_COUNT).fill(null)

let activeAnim = null
let isInitRollDone = false
let lastRollResult = null

const loadTgsData = async (index) => {
    if (animData[index]) return animData[index]

    try {
        const resp = await fetch(`${TGS_DIR_PREFIX}/dice_${index}.tgs`)

        if (!resp.ok) throw new Error(`HTTP ${resp.status}`)

        const buf = await resp.arrayBuffer()
        const raw = pako.inflate(new Uint8Array(buf), { to: "string" })
        animData[index] = JSON.parse(raw)

        return animData[index]
    } catch (err) {
        console.error(`Failed to load dice ${index}:`, err)

        return null
    }
}

const init = async () => {
    await Promise.all(Array.from({ length: DIE_ANIMATION_COUNT }, (_, i) => loadTgsData(i + 1)))
}

const getRandomLoaded = () => {
    const loaded = animData.reduce((acc, data, i) => {
        if (data && i !== lastRollResult) acc.push(i)
        return acc
    }, [])

    if (loaded.length === 0) {
        return Math.floor(Math.random() * DIE_ANIMATION_COUNT)
    }

    return loaded[Math.floor(Math.random() * loaded.length)]
}

const executeRoll = (value) => {
    const el = document.getElementById("die")
    el.classList.add("noclick")

    if (activeAnim) activeAnim.destroy()
    activeAnim = null

    const data = animData[value]
    if (!data) {
        el.classList.remove("noclick")
        return
    }

    activeAnim = lottie.loadAnimation({
        container: el,
        animationData: data,
        loop: false,
        autoplay: true
    })

    lastRollResult = value

    if (isInitRollDone) {
        document.getElementById("result").classList.add("opacity-0")
    }

    activeAnim.addEventListener("complete", () => {
        if (isInitRollDone) {
            const result = ROLL_RESULTS[value - 1]
            const resultEl = document.getElementById("result")
            resultEl.classList.remove("opacity-0")
            const img = `<img src="${result.icon}" alt="${result.alt}" class="h-6 w-6">`
            if (result.link) {
                resultEl.innerHTML = `<a href="${result.link}" target="_blank" class="inline-flex items-center gap-2">${img}${result.text}</a>`
            } else if (result.text) {
                resultEl.innerHTML = `<span class="inline-flex items-center gap-2">${img}${result.text}</span>`
            } else {
                resultEl.innerHTML = img
            }
        }

        isInitRollDone = true
        el.classList.remove("noclick")
    })
}

document.getElementById("die").addEventListener("click", () => executeRoll(getRandomLoaded()))
init().then(() => executeRoll(1))
