let copyOutputElem = document.querySelector("#copy-output");
let copyIndicatorElem = document.querySelector("#copy-indicator");
let outputElem = document.querySelector("#output");

function htmlDecode(input) {
    let doc = new DOMParser().parseFromString(input, "text/html");
    return doc.documentElement.textContent;
}

function showCopyIndicator(message) {
    copyIndicatorElem.innerText = message;
    copyIndicatorElem.classList.remove("hidden");
}

if (!!navigator.clipboard) {
    copyOutputElem.addEventListener("click", async () => {
        let outputInnerText = outputElem.innerHTML;

        try {
            await navigator.clipboard.writeText(htmlDecode(outputInnerText));
            showCopyIndicator("Copied");
        } catch (e) {
            console.log(e);
            if (e instanceof DOMException) {
                showCopyIndicator(`Clipboard access not available: ${e.name}`);
            } else {
                showCopyIndicator("Unknown error");
            }
        }
    });
    copyOutputElem.classList.remove("hidden");
}