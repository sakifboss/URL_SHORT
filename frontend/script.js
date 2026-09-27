const shortenForm = document.querySelector("#shorten-form");
const urlInput = document.querySelector("#long-url");
const submitButton = shortenForm.querySelector('button[type="submit"]');
const resultSection = document.querySelector("#result");
const shortUrlOutput = document.querySelector("#short-url");
const copyButton = document.querySelector("#copy-button");
const copyStatus = document.querySelector("#copy-status");
const errorMessage = document.querySelector("#error-message");

shortenForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    resultSection.hidden = true;
    errorMessage.hidden = true;
    copyStatus.textContent = "";
    submitButton.disabled = true;
    submitButton.textContent = "Shortening...";

    try {
        const response = await fetch("/shorten", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ url: urlInput.value.trim() }),
        });
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Unable to shorten this URL.");
        }

        if (typeof data.short_code !== "string") {
            throw new Error("The server returned an invalid response.");
        }

        shortUrlOutput.textContent = new URL(
            `/short/${encodeURIComponent(data.short_code)}`,
            window.location.origin,
        ).href;
        resultSection.hidden = false;
    } catch (error) {
        errorMessage.textContent = error instanceof Error
            ? error.message
            : "Unable to connect to the URL shortener.";
        errorMessage.hidden = false;
    } finally {
        submitButton.disabled = false;
        submitButton.textContent = "Shorten URL";
    }
});

copyButton.addEventListener("click", async () => {
    try {
        await navigator.clipboard.writeText(shortUrlOutput.textContent);
        copyStatus.textContent = "Copied to clipboard.";
    } catch {
        copyStatus.textContent = "Copy failed. Select the URL above and copy it manually.";
    }
});