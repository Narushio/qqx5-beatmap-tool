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
        formData.append('bpm', bpmInput.val());
        formData.append("file", fileInput[0].files[0]);


        $("#submit").addClass("is-loading")
        $.ajax({
            url: "/api/v1/update_xml_bpm",
            data: formData,
            type: "POST",
            contentType: false,
            processData: false,
            success: function (data) {
                $(".message-header").text("替换<Normal></Normal>标签")
                $("#output-box").removeClass("is-hidden")
                $("#submit").removeClass("is-loading")
                $("#notes").html(data)
            },
            error: function (data) {
                alert(data)
            },
        })
    })
}

initEvents()