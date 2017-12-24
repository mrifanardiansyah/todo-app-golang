var xhttp = new XMLHttpRequest();

function deleteItem(id) {
    $.ajax({
        url: '/api/todo/' + id,
        type: 'DELETE',
        success: function(result) {
            $('#item-' + id).remove();
        }
    })
}

function editItem(id) {
    window.location.href = "/todo/edit/" + id;
}

function gotoDetail(id) {
    window.location.href = "/todo/" + id;
}