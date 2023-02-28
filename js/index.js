function OSUFile(file = null, key = null, songTitle = null, songArtist = null, trackNotes = null, fileContent = null) {
	this.__proto__.file = file
	this.__proto__.key = key
	this.__proto__.songTitle = songTitle
	this.__proto__.songArtist = songArtist
	this.__proto__.trackNotes = trackNotes
	this.__proto__.fileContent = fileContent
}

function OSUNote(osuTrack = null, osuTime = null, osuNoteType = null, osuEndTime = null) {
	this.__proto__.osuTrack = osuTrack
	this.__proto__.osuTime = osuTime
	this.__proto__.osuNoteType = osuNoteType
	this.__proto__.osuEndTime = osuEndTime
}

function XMLNote(xmlTrack = null, xmlNoteType = null, bar = null, pos = null, endBar = null, endPos = null) {
	this.__proto__.xmlTrack = xmlTrack
	this.__proto__.xmlNoteType = xmlNoteType
	this.__proto__.bar = bar
	this.__proto__.pos = pos
	this.__proto__.endBar = endBar
	this.__proto__.endPos = endPos
}

const hitObject = {
	FOURKEY: {
		TRACK_ONE: "64",
		TRACK_TWO: "192",
		TRACK_THREE: "320",
		TRACK_FOUR: "448",
		LONG_NOTE: "128",
		SHORT_NOTE: ["1", "5"]
	},
	FIVEKEY: {
		TRACK_ONE: "51",
		TRACK_TWO: "153",
		TRACK_THREE: "256",
		TRACK_FOUR: "358",
		TRACK_FIVE: "460",
		LONG_NOTE: "128",
		SHORT_NOTE: ["1", "5"]
	}
}

let {
	file,
	key,
	songTitle,
	songArtist,
	trackNotes,
	fileContent
} = new OSUFile().__proto__

const initEvents = () => {
	$(".file-input").change((event) => {
		file = event.currentTarget.files[0]
		if (file !== undefined) {
			$(".file-name").html(file.name)
			loadFile()
		} else {
			$(".file-name").html("")
		}
	})

	$(".button.is-info").click(() => {
		if ($(".file-input").val() == "") {
			alert("请先选择文件呀！！！")
			return false
		}
		if ($(".input.is-info").val() == "") {
			alert("请输入BPM呀！！！")
			return false
		}
		$(".button.is-info").addClass("is-loading")
		$(".message-header").html(songArtist + " - " + songTitle)
		$("#output-box").removeClass("is-hidden")
		convertTimeline(key)
	})
}

const loadFile = () => {
	const reader = new FileReader()
	reader.readAsText(file, "UTF-8")
	reader.onload = (event) => {
		fileContent = event.target.result.replace(/\r\n/g, "newline ")
		key = fileContent.split("CircleSize:")[1].split(/newline/g)[0]
		songTitle = fileContent.split("Title:")[1].split(/newline/g)[0]
		songArtist = fileContent.split("Artist:")[1].split(/newline/g)[0]
		trackNotes = fileContent.split("[HitObjects]newline ")[1].split("newline ")
		trackNotes.pop()
	}
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
	const bpm = $(".input.is-info").val()
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
			if ((Number(osuTime) % barTime / posTime) > (Math.round(Number(osuTime) % barTime /
					posTime))) {
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
					xmlTrack = "Left1
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
	$(".button.is-info").removeClass("is-loading")
}

initEvents()
