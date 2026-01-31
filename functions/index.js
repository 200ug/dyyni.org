export async function onRequest(context) {
    const ua = context.request.headers.get("user-agent") || ""
    if (!/w3m|lynx|links|elinks|curl|wget/.test(ua.toLowerCase())) {
      return context.next();
    }
    const resp = await context.next()
    return new HTMLRewriter()
        .on("header img", {
            element(el) {
                el.replace(`<pre>

   ___/ /_ ____ _____  (_)___  _______ _
  / _  / // / // / _ \\/ // _ \\/ __/ _ \`/
  \\_,_/\\_, /\\_, /_//_/_(_)___/_/  \\_, / 
      /___//___/                 /___/  

</pre>`, {html: true})
            }
    }).transform(resp)
}
