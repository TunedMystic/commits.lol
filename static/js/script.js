const header = `
  _   _ _   _____ _                   _
 | | | (_) |_   _| |__   ___ _ __ ___| |
 | |_| | |   | | | '_ \\ / _ \\ '__/ _ \\ |
 |  _  | |   | | | | | |  __/ | |  __/_|
 |_| |_|_|   |_| |_| |_|\\___|_|  \\___(_)
\n\n`;

document.addEventListener('DOMContentLoaded', (e) => {
    console.log(header);
    document.getElementById('censor-toggle').addEventListener('click', toggleCensoredCommits);
    document.getElementById('more-commits').addEventListener('click', fetchMoreCommits());
    document.getElementById('poop-commits').addEventListener('click', fetchMoreCommits('poop'));
    document.getElementById('holy-commits').addEventListener('click', fetchMoreCommits('holy'));
});

// This function shows / hides the bad words in every commit message.
// It also swaps the emoji based on the state of the censored message.
function toggleCensoredCommits() {
    const curseFace = '/static/img/curse-face.png';
    const flushedFace = '/static/img/flushed-face.png';
    const img = document.querySelector('#censor-toggle > img');

    // Swap emojis.
    if (img.getAttribute("src") === curseFace) {
        img.setAttribute("src", flushedFace);

    } else if (img.getAttribute("src") === flushedFace) {
        img.setAttribute("src", curseFace);
    }

    // Swap the commit message by toggling a class on the parent element.
    document.getElementById('commit-items').classList.toggle('uncensor');
}

// This function fetches a new batch of random commits and places it in the commits container.
function fetchMoreCommits(group) {
    group = typeof(group) === 'undefined' ? '' : group;
    return function() {
        makeRequest(`/?group=${group}&fragment=true`, (html) => {
            document.getElementById('commit-items').innerHTML = html;
        }, 'when fetching commits');
    }
}

// ------------------------------------------------------------------
// Utils / Helper functions
// ------------------------------------------------------------------
const errMessage = (reason) => `Something went wrong ${reason}.`;

const makeRequest = (url, onSuccess, errorMessage) => {
    fetch(url)
    .then((response) => {
        if (!response.ok) { throw new Error(); }
        return response.text();
    })
    .then(html => onSuccess(html))
    .catch((err) => {
        // If a custom message was caught, then use it. Else, construct the error message with the given reason.
        const errorMessageText = err.message ? err.message : errMessage(errorMessage);
        alert(errorMessageText);
    });
};
