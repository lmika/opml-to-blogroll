let copyOutputElem = document.querySelector("#copy-output");
let copyIndicatorElem = document.querySelector("#copy-indicator");
let outputElem = document.querySelector("#output");

function showCopyIndicator(message) {
    copyIndicatorElem.innerText = message;
    copyIndicatorElem.classList.remove("hidden");
}

if (!!navigator.clipboard) {
    copyOutputElem.addEventListener("click", async () => {
        let outputInnerText = outputElem.innerHTML;

        try {
            await navigator.clipboard.write(outputInnerText);
            showCopyIndicator("Copied");
        } catch (e) {
            console.log(e);
            if (e instanceof NotAllowedError) {
                showCopyIndicator("Page not permitted to use clipboard");
            } else {
                showCopyIndicator("Unknown error");
            }
        }
    });
    copyOutputElem.classList.remove("hidden");
}