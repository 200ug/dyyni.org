const ENDPOINT = "https://organic.dyyni.org/blackbox"
const MAX_MESSAGE_LENGTH = 80

const form = document.getElementById("blackbox-form")
const input = document.getElementById("blackbox-input")

input.maxLength = MAX_MESSAGE_LENGTH

// TODO: add a server response handle (simply, change the input field border color to red on error)

form.addEventListener("submit", (e) => {
    e.preventDefault()

    const message = input.value.trim()
    if (!message) return

    fetch(ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message })
    })

    input.value = ""
})

