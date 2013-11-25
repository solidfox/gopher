<?php

	$query = $_GET["query"];


	$data = array(


		array(
			'score' => 12,
			'title' => 'hihi',
			'url' => 'http://lalala.com',
			'modificationDate' => '10000000',
			'size' => 12,
			'keywords' => array(

				array('word'=> 'test',	'freq'=> 3),
				array('word'=> 'baby',	'freq'=> 0)

				),
			'parent' => array(

				array('url' =>'http://wht.com', 'title' =>'test'),
				array('url' =>'http://wht.com', 'title' =>'test'),
				),

			'child' => array(
				array('url' =>'http://wht.com', 'title' =>'test'),
				array('url' =>'http://wht.com', 'title' =>'test'),
				)


			),


		array(
			'score' => 10,
			'title' => 'scrum',
			'url' => 'http://Scrum.com',
			'modificationDate' => '10000000',
			'size' => 13,
			'keywords' => array(

				array('word'=> 'test',	'freq'=> 1),
				array('word'=> 'baby',	'freq'=> 1)

				),

			'parent' => array(

				array('url' =>'http://wht.com', 'title' =>'test'),
				array('url' =>'http://wht.com', 'title' =>'test'),
				),

			'child' => array(
				array('url' =>'http://wht.com', 'title' =>'test'),
				array('url' =>'http://wht.com', 'title' =>'test'),
				)

			)

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