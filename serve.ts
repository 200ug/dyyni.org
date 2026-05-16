const STATIC_DIR = "dist"

Bun.serve({
    port: 3000,
    async fetch(req) {
        const url = new URL(req.url)
        const path =
            url.pathname === "/" ? `${STATIC_DIR}/index.html` : `${STATIC_DIR}${url.pathname}`
        const file = Bun.file(path)
        if (await file.exists()) return new Response(file)
        return new Response("Not found", { status: 404 })
    }
})

console.log("dev server running: http://localhost:3000")
