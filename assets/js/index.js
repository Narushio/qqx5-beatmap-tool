const initEvents = () => {
	let fileContent, songTitle, songArtist

	$(".file-input").change((event) => {
		let file = event.currentTarget.files[0]

		if (file !== undefined) {
			$(".file-name").html(file.name)
			const reader = new FileReader()
			reader.readAsText(file, "UTF-8")
			reader.onload = (event) => {
				fileContent = event.target.result.replace(/\r\n/g, "newline ")
				songTitle = fileContent.split("Title:")[1].split(/newline/g)[0]
				songArtist = fileContent.split("Artist:")[1].split(/newline/g)[0]
			}
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
			url: "/api/v1/convert/osu_to_xml",
			data: formData,
			type: "POST",
			contentType: false,
			processData: false,
			success: function (data) {
				$("#creator").text("Creator: " + data.creator)
				$("#offset").text("Offset: " + data.offset)
				$("#output-box-title").text(songTitle + " - " + songArtist)
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
