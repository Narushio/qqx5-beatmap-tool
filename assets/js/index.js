const initEvents = () => {
	$(".file-input").change((event) => {
		let file = event.currentTarget.files[0]
		let fileContent, songTitle, songArtist

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
		if (fileInput.val() === "") {
			alert("Please select a file!")
			return
		}
		const formData = new FormData();
		formData.append("file", fileInput[0].files[0]);


		$("#submit").addClass("is-loading")
		$.ajax({
			url: "/api/v1/convert/osu_to_xml",
			data: formData,
			type: "POST",
			contentType: false,
			processData: false,
			success: function (data) {
				$("#output-box-title").text(songTitle + " - " + songArtist)
				$("#output-box").removeClass("is-hidden")
				$("#submit").removeClass("is-loading")
				$("#notes").html(data)
			},
			error: function (data) {
				alert(data.responseText)
				$("#submit").removeClass("is-loading")
			},
		})
	})


	// $("#submit").click(() => {
	// 	if ($(".file-input").val() == "") {
	// 		alert("Please select a file!")
	// 		return false
	// 	}
	// 	if ($(".input").val() == "") {
	// 		alert("Please enter BPM!")
	// 		return false
	// 	}
	//
	// 	$("#submit").addClass("is-loading")
	// 	$(".message-header").html(songArtist + " - " + songTitle)
	// 	$("#output-box").removeClass("is-hidden")
	// 	convertTimeline(key)
	// })
}

const noteTempStr = (bar, pos, track, type, endBar, endPos) => {
	let note = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&lt;Note Bar=\"" + bar + "\" Pos=\"" + pos +
		"\" from_track=\"" + track + "\" target_track=\"" + track + "\" note_type=\"" + type + "\" ";
	if (endBar != null) {
		note += "EndBar=\"" + endBar + "\" EndPos=\"" + endPos + "\"";
	}
	note += " /&gt;";
	return note;
}

const convertTimeline = (key) => {
	const bpm = $(".input").val()
	const barTime = 60000 / bpm * 4
	const posTime = barTime / 32
	let xmlContent = ""

	let {
		osuTrack,
		osuTime,
		osuNoteType,
		osuEndTime
	} = new OSUNote().__proto__

	let {
		xmlTrack,
		xmlNoteType,
		bar,
		pos,
		endBar,
		endPos,
	} = new XMLNote().__proto__


	trackNotes.forEach((note, index) => {
		osuTrack = note.split(",")[0]
		osuTime = note.split(",")[2]
		osuNoteType = note.split(",")[3]

		bar = Math.floor(Number(osuTime) / barTime) + 1
		pos = Math.round(Number(osuTime) % barTime / posTime) * 2
		
		if ((pos % 2) !== 0) {
			if ((Number(osuTime) % barTime / posTime) > (Math.round(Number(osuTime) % barTime / posTime))) {
				pos += 1
			} else {
				pos -= 1
			}
		}

		if (osuNoteType === hitObject.FOURKEY.LONG_NOTE || osuNoteType === hitObject.FIVEKEY.LONG_NOTE) {
			xmlNoteType = "long"
			osuEndTime = note.split(",")[5].split(":")[0]
			endBar = Math.floor(Number(osuEndTime) / barTime) + 1
			endPos = Math.round(Number(osuEndTime) % barTime / posTime) * 2
			if ((endPos % 2) !== 0) {
				if ((Number(osuEndTime) % barTime / posTime) > (Math.round(Number(osuEndTime) %
						barTime / posTime))) {
					endPos += 1
				} else {
					endPos -= 1
				}
			}
		} else {
			endBar = undefined
			endPos = undefined
			xmlNoteType = "short"
		}
		
		if (key === "4") {
			switch (osuTrack) {
				case hitObject.FOURKEY.TRACK_ONE:
					xmlTrack = "Left2"
					break
				case hitObject.FOURKEY.TRACK_TWO:
					xmlTrack = "Left1"
					break
				case hitObject.FOURKEY.TRACK_THREE:
					xmlTrack = "Right1"
					break
				case hitObject.FOURKEY.TRACK_FOUR:
					xmlTrack = "Right2"
					break
			}
		}
		
		if (key === "5") {
			switch (osuTrack) {
				case hitObject.FIVEKEY.TRACK_ONE:
					xmlTrack = "Left2"
					break
				case hitObject.FIVEKEY.TRACK_TWO:
					xmlTrack = "Left1"
					break
				case hitObject.FIVEKEY.TRACK_THREE:
					xmlTrack = "Middle"
					break
				case hitObject.FIVEKEY.TRACK_FOUR:
					xmlTrack = "Right1"
					break
				case hitObject.FIVEKEY.TRACK_FIVE:
					xmlTrack = "Right2"
					break
			}
		}
		
		xmlContent += noteTempStr(bar, pos, xmlTrack, xmlNoteType, endBar, endPos) + "<br/>"
	})
	$("#resultDiv").html(xmlContent)
	$("#submit").removeClass("is-loading")
}

initEvents()
