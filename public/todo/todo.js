(function ($) {
    'use strict';

    $(function () {
        var filter = 'all';
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');

        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                $.post("/todos", { name: item }, addItem)
                todoListInput.val("");
            }
        });

        var addItem = function (item) {
            var completedClass = item.completed ? "completed" : "";
            var picture = item.picture ? item.picture : "";
            var createdAt = new Date(item.created_at);
            createdAt = createdAt.getFullYear() + '年' +
                ('0' + (createdAt.getMonth() + 1)).slice(-2) + '月' +
                ('0' + createdAt.getDate()).slice(-2) + ' ' +
                ('0' + createdAt.getHours()).slice(-2) + '時' +
                ('0' + createdAt.getMinutes()).slice(-2)+ '分';
            var listItemHtml =
                "<li class='" +
                completedClass +
                "' id='" +
                item.id +
                "' style='display: flex; justify-content: space-between;'>" +
                "<div class='profile-image-container' style='margin-right: 10px;'>" +
                "<img class='profile-image' src='" + picture + "' style='width:25px; height:25px;'/>" +
                "</div>" +
                "<div style='align-self: center; flex: 1;'>" +
                "<div class='form-check'>" +
                "<label class='form-check-label'>" +
                "<input class='checkbox' type='checkbox'" +
                (item.completed ? " checked='checked'" : "") +
                "/>" +
                item.name +
                "<i class='input-helper'></i></label></div></div>" +
                "<span class='createdAt' style='margin-right: 10px; color: #999;'>" + createdAt + "</span>" +
                "<i class='remove mdi mdi-close-circle-outline'></i>" +
                "</li>";
            todoListItem.append(listItemHtml);

            updateProgressBar();
        };


        $.get('/todos', function (items) {
            items.forEach(e => {
                addItem(e)
            });
        });

        todoListItem.on('change', '.checkbox', function () {
            var id = $(this).closest("li").attr('id')
            var $self = $(this);
            var complete = true;
            if ($(this).attr('checked')) {
                complete = false;
            }
            $.get("/complete-todo/" + id + "?complete=" + complete, function (data) {
                if (complete) {
                    $self.attr('checked', 'checked');
                } else {
                    $self.removeAttr('checked');
                }
                $self.closest("li").toggleClass('completed');

                updateProgressBar();
            })
        });

        todoListItem.on('click', '.remove', function () {
            var id = $(this).closest("li").attr('id');
            var $self = $(this);
            $.ajax({
                url: "/todos/" + id,
                type: "DELETE",
                success: function (data) {
                    if (data.success) {
                        $self.parent().remove();
                        updateProgressBar();
                    }
                }
            })
        });

        $('.filter-btn').click(function () {
            $(this).addClass('active').siblings().removeClass('active');

            filter = $(this).data('filter');
            if (filter === 'user') {
                sortByUser();
            } else if (filter === 'completed') {
                sortByCompleted();
            } else {
                showAll();
            }
        });

        $('.completed-clear-btn').click(function () {
            $.ajax({
                url: "/todos-completed-clear",
                type: "DELETE",
                success: function (data) {
                    if (data.success) {
                        $self.parent().remove();
                        updateProgressBar();
                        location.reload();
                    }
                }
            })
        });

        const upArrowBtn = document.querySelector('.up-arrow');
        const downArrowBtn = document.querySelector('.down-arrow');

        upArrowBtn.addEventListener('click', () => {
            upArrowBtn.style.display = 'none';
            downArrowBtn.style.display = 'block';
            var sort = 'desc';
            $.get('/todos/sorted', { filter: filter, sort: sort }, function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        });

        downArrowBtn.addEventListener('click', () => {
            downArrowBtn.style.display = 'none';
            upArrowBtn.style.display = 'block';
            var sort = 'asc';
            $.get('/todos/sorted', { filter: filter, sort: sort }, function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        });

        function sortByUser() {
            $.get('/todos/sorted-by-user', function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        }

        function sortByCompleted() {
            $.get('/todos/sorted-by-completed', function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        }

        function showAll() {
            $.get('/todos', function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        }

        function clearList() {
            $('.todo-list').empty();
        }

        function updateProgressBar() {
            $.get('/todos/progress', function (progress) {
                $('.progress-bar').css('width', progress + '%');
                $('.progress-bar').text(progress + '%');
            });
        }
    });
})(jQuery);