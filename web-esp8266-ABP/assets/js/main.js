$(function() {

    setInterval(function() {
        getDisplays()
    }, 2000);

});

function setDisplays(id, dispositivosClient) {

    $(`#sensor-client-${id}`).html(dispositivosClient.sensor + ' cm')

    $(`#sala-client-${id}`).val(dispositivosClient.sala)
    $(`#quarto-client-${id}`).val(dispositivosClient.quarto)
    $(`#cozinha-client-${id}`).val(dispositivosClient.cozinha)

    let checkboxAll = document.querySelectorAll(`.client-${id}`);
    let statusLight = ''

    for (var i = 0; i < checkboxAll.length; i++) {
        // utilize o valor aqui, adicionei ao array para exemplo
        let checkboxItem = checkboxAll[i]
        statusLight = $(`#light-${checkboxItem.name}`)

        if (checkboxItem.value == 1) {
            statusLight.html(`<img src="./assets/img/light-bulb-on.png" alt="" style="height:65px;">`)
            checkboxItem.className += ' switch-shadow-special'
            checkboxItem.checked = true

        } else {
            statusLight.html(`<img src="./assets/img/turned-off.png" alt="" style="height:65px;">`)
            checkboxItem.checked = false
        }
    }
}


$(".switch").click(function(event) {

    let checkboxItem = event.target
    let checkboxId = checkboxItem.id

    let statusLight = $(`#light-${checkboxItem.name}`)

    if (checkboxItem.checked == true) {
        checkboxItem.value = 1
        statusLight.html(`<img src="./assets/img/light-bulb-on.png" alt="" style="height:65px;">`)
        checkboxItem.className += ' switch-shadow-special'

    } else {
        checkboxItem.value = 0
        statusLight.html(`<img src="./assets/img/turned-off.png" alt="" style="height:65px;">`)
        checkboxItem.classList.remove('switch-shadow-special')
    }

    postDisplaysLight(checkboxId, checkboxItem.value)

});

function postDisplaysLight(param, status) {

    const [local, user, id] = param.split("-")
        /* console.log(local)
        console.log(user)
        console.log(id) */

    let url = `http://localhost:8000/${local}/${id}/${status}`

    console.log(url)
    $.ajax(url, // request url
        {
            success: function(data, status, xhr) { // success callback function
                console.log(data)
            }
        });


}

function getDisplays() {

    /* 
    //* EXEMPLO DE RETORNO

    let returnAjax = [{
            "clientId": "1",
            "sala": "on",
            "quarto": "on",
            "cozinha": "off",
            "sensor": "25 cm"
        }, {
            "clientId": "2",
            "sala": "on",
            "quarto": "on",
            "cozinha": "off",
            "sensor": "25 cm"
        }] */

    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            console.log('carregou')
        }
    };
    xhttp.open("GET", "http://localhost:8000/", true);
    xhttp.onload = function(e) {
        var arraybuffer = JSON.parse(xhttp.response); // não é responseText
        /* ... */
        arraybuffer.map(function(nome, i) {
            //console.log('[forEach]', nome, i);
            setDisplays(nome.clientId, nome)

        })
    }
    xhttp.send();

}