$(document).ready(function () {


	var currentUrl = window.location.href
	console.log(currentUrl)
	//console.log(currentUrl.replace("login", "signup"))
		//console.log(currentUrl.replace("signup", ""))
		//console.log(String(currentUrl) - "/signup")

	$("#loginbutton").click(function () {
		console.log($("#userName").val())
		console.log($("#password").val())


		window.userName = $("#userName").val();
		window.password = $("#password").val();


		$.ajax({
			data: {
				"username": window.userName,
				"password": window.password,

			},
			//dataType: "json",
			type: "POST",

			url: currentUrl + "submit",


			success: function (data) {
				//$('#loginresult').empty()
				//console.log("Data sent: ", data)
				$('#loginresult').append(data);

			},
			error: function (req, status, err) {
				//console.log(req.responseText)
				console.log(req)

				console.log('Something went wrong', status, err);
				console.log(err)

			}
		});


	});
	
	
	$("#trialscrape").click(function () {
		
		$.ajax({
			data: {
			
			},
			dataType: "json",
			type: "POST",

			url: currentUrl.replace("login", "trialscrapelocation"),
			
		

			success: function (data) {
				$("#ajaxResults").empty();
				console.log(data)

				console.log(data[1].Title)
				console.log(data.Cost)
				for (var i = 0; i < data.length; i++) {
					var imagetag = "<img src=" + data[i].ImageUrl + " style='width:100px;height:100px;'>";

					$('#ajaxResults').append(data[i].Title + " " + data[i].Cost + "<br>" + imagetag + "<br><br>");


				}

			},
			error: function (req, status, err) {
				$("#ajaxResults").empty();
				console.log(req.responseText)
				$('#ajaxResults').append(req.responseText);
				console.log('Something went wrong', status, err);
				console.log(err)

			}
		});
	});
	
	
	$("#signupredirect").click(function (e) {
		e.preventDefault();
		console.log(currentUrl.replace("login", "signup"))
		window.location = currentUrl.replace("login", "signup");
	});
	$("#homepageredirect").click(function (e) {
		e.preventDefault();
		//console.log(currentUrl.replace("login", ""))
		window.location = currentUrl.replace("login", "");
	});
	$("#logout").click(function (e) {
		
		e.preventDefault();
		//console.log(currentUrl.replace("login", ""))
		window.location = currentUrl.replace("login", "logout");
	});

});