<?php

	$query = $_GET["url"];


	$data = array(
		'scrum',
		'essay',
		'gohper',
		'southpark'
	 );


// score
// page title
// url
// last modification date, size of page
// keyword 1 freq 1; keyword 2 freq 2; …
// Parent link 1
// Parent link 2
// … …
// Child link 1
// Child link 2


	// Send the data.
	echo json_encode($data);




?>