meetingGuid = window.location.href.split('?')[1];
meetingURL = 'https://diarizer.blabbertabber.com/' + meetingGuid + '/diarization.txt';

jQuery('#xyz').html("jQuery works.  But we haven't downloaded the file yet");

jQuery.get(meetingURL, function(data) {
  alert("We are within the get function.")
  jQuery('#xyz').html("we hit the get file function!!!!");
});

jQuery('#last').html("really done" + new Date());
