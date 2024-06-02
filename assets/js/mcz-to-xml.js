const initEvents = () => {
    $(".file-input").change((event) => {
        file = event.currentTarget.files[0]
        if (file !== undefined) {
            $(".file-name").html(file.name)
        } else {
            $(".file-name").html("")
        }
    })

    $("#submit").click(() => {
        const fileInput = $(".file-input")
        const bpmInput = $(".input")
        if (fileInput.val() === "") {
            alert("Please select a file!")
            return
        }
        if (bpmInput.val() === "") {
            alert("Please enter BPM!")
            return
        }

        const formData = new FormData();
        formData.append("bpm", bpmInput.val());
        formData.append("file", fileInput[0].files[0]);


        $("#submit").addClass("is-loading")
        $.ajax({
            url: "/api/v1/convert/mcz_to_xml",
            data: formData,
            type: "POST",
            contentType: false,
            processData: false,
            success: function (data) {
                $("#creator").text("Creator: " + data.creator)
                $("#offset").text("Offset: " + data.offset)
                $("#output-box-title").text(data.title + " - " + data.artist)
                $("#output-box").removeClass("is-hidden")
                $("#submit").removeClass("is-loading")
                $("#notes").html(data.notes)
            },
            error: function (data) {
                alert(data.responseText)
                $("#submit").removeClass("is-loading")
            },
        })
    })
}

initEvents()