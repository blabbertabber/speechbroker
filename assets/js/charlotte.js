// https://diarizer.blabbertabber.com?meeting=1766e8dc-28e1-11e7-a2c1-000c291285ff

meetingGuid = window.location.href.split('?')[1];
meetingGuid = meetingGuid.split('=')[1];
meetingURL = 'https://diarizer.blabbertabber.com/' + meetingGuid + '/diarization.txt';

jQuery('#xyz').html("jQuery works.  But we haven't downloaded the file yet");

jQuery.get(meetingURL, function (data) {
    var speakerTimes = {};
    var lines = data.split(/\n/);
    var out = "";
    lines.forEach(function (line) {
        if (line.length > 0) {
            var fields = line.split(/\s+/);
            var startTime = fields[2].split(/=/)[1];
            var endTime = fields[3].split(/=/)[1];
            var speakerNum = fields[4].split(/=/)[1].split(/_/)[1];
            out = out + startTime + " " + endTime + " " + speakerNum + "<br />";
            if (!speakerTimes[speakerNum] || !("time" in speakerTimes[speakerNum])) {
                speakerTimes[speakerNum] = {time: endTime - startTime};
                speakerTimes[speakerNum].time = 0
            }
            speakerTimes[speakerNum].time +=  endTime - startTime;
        }
    });
    jQuery('#asdf').html(out); // + new Date());
    jQuery('#xyz').html(speakerTimes["1"].time);

});

