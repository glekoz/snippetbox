window.onload = _ => {
    [...document.getElementsByClassName('content')].forEach(e => e.innerText = e.innerText.replaceAll("\\n", '\n'));
}