import lottie from "lottie-web"
import pako from "pako"

const TGS_DIR_PREFIX = "/static/dice"
const DIE_ANIMATION_COUNT = 7

// null (not loaded) | object (decompressed lottie json)
const animData = new Array(DIE_ANIMATION_COUNT).fill(null)

let activeAnim = null
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
    await Promise.all(Array.from({ length: DIE_ANIMATION_COUNT }, (_, i) => loadTgsData(i)))
}

const getRandomLoaded = () => {
    const loaded = animData.reduce((acc, data, i) => {
        if (data && i !== 0 && i !== lastRollResult) acc.push(i)
        return acc
    }, [])

    if (loaded.length === 0) {
        return 1 + Math.floor(Math.random() * (DIE_ANIMATION_COUNT - 1))
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

    document.getElementById("result").textContent = `Rolled: ${value}`
    lastRollResult = value

    activeAnim.addEventListener("complete", () => el.classList.remove("noclick"))
}

document.getElementById("die").addEventListener("click", () => executeRoll(getRandomLoaded()))
init().then(() => executeRoll(1))
