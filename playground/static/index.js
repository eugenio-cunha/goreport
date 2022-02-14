'use strict';

const HtmlMode = ace.require("ace/mode/html").Mode;
const JsonMode = ace.require("ace/mode/json").Mode;

const header = ace.edit("header");
header.setTheme("ace/theme/monokai");
header.session.setMode(new HtmlMode());

const body = ace.edit("body");
body.setTheme("ace/theme/monokai");
body.session.setMode(new HtmlMode());

const footer = ace.edit("footer");
footer.setTheme("ace/theme/monokai");
footer.session.setMode(new HtmlMode());

const data = ace.edit("data");
data.setTheme("ace/theme/monokai");
data.session.setMode(new JsonMode());

document.getElementById("defaultOpen").click();

function openEditor(evt, name) {
  // Declare all variables
  var i, tabcontent, tablinks;

  // Get all elements with class="tabcontent" and hide them
  tabcontent = document.getElementsByClassName("tabcontent");
  for (i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }

  // Get all elements with class="tablinks" and remove the class "active"
  tablinks = document.getElementsByClassName("tablinks");
  for (i = 0; i < tablinks.length; i++) {
    tablinks[i].className = tablinks[i].className.replace(" active", "");
  }

  // Show the current tab, and add an "active" class to the button that opened the tab
  document.getElementById(name).style.display = "block";
  evt.currentTarget.className += " active";
}

async function preview() {
  const payload = {
    header: header.getValue(),
    body: body.getValue(),
    footer: footer.getValue(),
    data: data.getValue(),
  }

  const res = await window.fetch('http://localhost:8080/report', {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(payload),
    cache: 'default'
  })
  const blob = await res.blob()
  let blobURL = URL.createObjectURL(blob);
  document.querySelector("iframe").src = blobURL;
}