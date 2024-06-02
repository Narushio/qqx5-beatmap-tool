let clipboard = new ClipboardJS("#btn-copy");
clipboard.on("success", function (e) {
    e.clearSelection();
    alert("Copied :)")
});

clipboard.on("error", function (e) {
    alert("Uncopied :( Please report developer this bug.")
});