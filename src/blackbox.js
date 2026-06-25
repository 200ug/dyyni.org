const ENDPOINT = ENDPOINT_URL
const MAX_MESSAGE_LENGTH = 160
const PLACEHOLDERS = [
    "What's up?",
    "How are you?",
    "Say hello...",
    "Leave a note...",
    "Write a poem...",
    "Drop an idea...",
]

const form = document.getElementById("blackbox-form")
const input = document.getElementById("blackbox-input")

input.placeholder = PLACEHOLDERS[Math.floor(Math.random() * PLACEHOLDERS.length)]
input.maxLength = MAX_MESSAGE_LENGTH

function flashBorder(color) {
    input.style.setProperty("--blink-color", color)
    input.classList.remove("blink-border")
    void input.offsetWidth
    input.classList.add("blink-border")
    input.addEventListener(
        "animationend",
        () => {
            input.classList.remove("blink-border")
            input.style.removeProperty("--blink-color")
        },
        { once: true }
    )
}

form.addEventListener("submit", (e) => {
    e.preventDefault()

    const message = input.value.trim()
    if (!message) return

    fetch(ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message })
    })
        .then((res) => {
            flashBorder(res.ok ? "var(--color-success)" : "var(--color-error)")
        })
        .catch(() => {
            flashBorder("var(--color-error)")
        })

    input.value = ""
})
