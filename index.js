// Import the content HTML.
import html from "./content.html";

// Reveal dependencies.
import "headjs/dist/1.0.0/head.js";

// Reveal proper.
import Reveal from "reveal.js";

// Stylesheet imports.
import "reveal.js/css/reveal.css";
import "reveal.js/css/theme/sky.css";
import "reveal.js/lib/css/zenburn.css";
import "./pres.less";

// Actually write the content HTML.
document.write(html);

Reveal.initialize({
  controls: (-1 == window.location.host.indexOf("localhost") && -1 == window.location.host.indexOf("0.0.0.0")),
  progress: false,
  history: true,
  center: true,
  width: 1920,
  height: 1080,

  transition: 'fade', // none/fade/slide/convex/concave/zoom

  // Optional reveal.js plugins
  dependencies: [
    { src: 'node_modules/reveal.js/lib/js/classList.js', condition: function() { return !document.body.classList; } },
    { src: 'node_modules/reveal.js/plugin/highlight/highlight.js', async: true, condition: function() { return !!document.querySelector( 'pre code' ); }, callback: function() { hljs.initHighlightingOnLoad(); } },
    { src: 'node_modules/reveal.js/plugin/notes/notes.js', async: true },
  ],
});

// Make the Reveal object available to plugins.
window.Reveal = Reveal;
