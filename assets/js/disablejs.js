function favicons() {
  this.hidden = "hidden"
  this.visibilityChange = "visibilitychange"
  this.favicon = document.querySelector("[rel='shortcut icon']").href
  this.title = document.title
  this.wasSpoofed = false
  this.spoofed = []

  this.services = {
    sy: () => {
      const title = "Official Church of Scientology: Difficulties on the Job - Online Course"
      const favicon = "/images/disablejs/sy.png"
      return {
        title,
        favicon,
      }
    },
    cdc: () => {
      const title = "Ask HN: How could I safely contact drug cartels?"
      const favicon = "/images/disablejs/hn.png"
      return {
        title,
        favicon,
      }
    },
    rps: () => {
      const title = "rust programming socks - Google Shopping"
      const favicon = "/images/disablejs/go.ico"
      return {
        title,
        favicon,
      }
    },
    aec: () => {
      const title = "Adult entertainment clubs - Google Maps"
      const favicon = "/images/disablejs/gm.png"
      return {
        title,
        favicon,
      }
    },
    pul: () => {
      const title = "Pick up lines suggestions - ChatGPT"
      const favicon = "/images/disablejs/cgpt.png"
      return {
        title,
        favicon,
      }
    },
    fes: () => {
      const title = "The Flat Earth Society"
      const favicon = "/images/disablejs/fes.png"
      return {
        title,
        favicon,
      }
    },
    tsm: () => {
      const title = "Amazon.com: taylor swift merch"
      const favicon = "/images/disablejs/az.ico"
      return {
        title,
        favicon,
      }
    },
    wp: () => {
      const title = "Amazon.com: waifu pillow"
      const favicon = "/images/disablejs/az.ico"
      return {
        title,
        favicon,
      }
    },
    rwsb: () => {
      const title = "r/wallstreetbets on Reddit"
      const favicon = "/images/disablejs/rd.png"
      return {
        title,
        favicon,
      }
    },
    iw: () => {
      const title = "Infowars: There's a War on For Your Mind!"
      const favicon = "/images/disablejs/iw.png"
      return {
        title,
        favicon,
      }
    },
    tac: () => {
      const title = "The Anarchist Cookbook by William Powell | Goodreads"
      const favicon = "/images/disablejs/gr.png"
      return {
        title,
        favicon,
      }
    },
    nxfs: () => {
      const title = "Fifty Shades of Grey | Netflix"
      const favicon = "/images/disablejs/nx.ico"
      return {
        title,
        favicon,
      }
    },
    jbn: () => {
      const title = "jeff bezos nudes - Google Image Search"
      const favicon = "/images/disablejs/go.ico"
      return {
        title,
        favicon,
      }
    },
    nggyu: () => {
      const title = "Rick Astley - Never Gonna Give You Up - YouTube"
      const favicon = "/images/disablejs/yt.ico"
      return {
        title,
        favicon,
      }
    },
    beast: () => {
      const title = "MrBeast en Español - YouTube"
      const favicon = "/images/disablejs/yt.ico"
      return {
        title,
        favicon,
      }
    },
    ftx: () => {
      const title = "FTX Cryptocurrency Exchange"
      const favicon = "/images/disablejs/ftx.png"
      return {
        title,
        favicon,
      }
    },
  }

  this.enabledServices = Object.keys(this.services)

  this.init = function () {
    if (typeof document.mozHidden !== "undefined") {
      this.hidden = "mozHidden"
      this.visibilityChange = "mozvisibilitychange"
    } else if (typeof document.msHidden !== "undefined") {
      this.hidden = "msHidden"
      this.visibilityChange = "msvisibilitychange"
    } else if (typeof document.webkitHidden !== "undefined") {
      this.hidden = "webkitHidden"
      this.visibilityChange = "webkitvisibilitychange"
    }

    document.addEventListener(this.visibilityChange, this.handler.bind(this), false)
  }

  this.restore = function () {
    const title = this.title
    const favicon = this.favicon

    this.update({
      title,
      favicon,
    }, false) // cache busting disabled
  }

  this.update = function (data, bustCache) {
    const newHref = bustCache ? data.favicon + "?v=" + Math.round(Math.random() * 10000000) : data.favicon
    const newLink = document.createElement("link")
    newLink.type = "image/x-icon"
    newLink.rel = "shortcut icon"
    newLink.href = newHref

    document
      .getElementsByTagName("head")[0]
      .querySelector("[rel='shortcut icon']")
      .remove();
    document.getElementsByTagName("head")[0].appendChild(newLink);
    document.title = data.title;

    /*
    if (bustCache) {
      console.log(`new title: '${document.title}'`)
      console.log(`new href: '${newHref}'`)
    } else {
      console.log("restoring old favicon")
    }
    */
  }

  this.spoof = function () {
    let i = 0

    if (this.spoofed.length === this.enabledServices.length) {
      this.spoofed.length = 0
    }

    for (let es = 0; es < this.enabledServices.length; es++) {
      i = Math.round(Math.random() * (this.enabledServices.length - 1))
      if (this.spoofed.includes(i) === false) {
        break
      }
    }

    this.spoofed.push(i)
    const service = this.enabledServices[i]

    if (service && this.services[service]) {
      this.update(this.services[service](), true) // cache busting enabled
    }

    this.wasSpoofed = true
  }

  this.handler = function () {
    if (document[this.hidden]) {
      this.spoof()
    } else {
      this.restore()
    }
  }

  this.init()
}

function closeDisablejsInfo() {
  document.getElementById("disablejs").style.display = "none"
}
