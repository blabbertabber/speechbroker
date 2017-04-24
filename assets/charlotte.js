meetingGuid = window.location.href.split('?')[1];
meetingGuid = meetingGuid.split('=')[1];
meetingURL = 'https://diarizer.blabbertabber.com/' + meetingGuid + '/diarization.txt';

jQuery('#xyz').html("jQuery works.  But we haven't downloaded the file yet");

jQuery.get(meetingURL, function(data) {
  jQuery('#xyz').html(data.replace("\n", "<br />"));
});

jQuery('#asdf').html("really done"); // + new Date());
