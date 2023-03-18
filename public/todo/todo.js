(function ($) {
    'use strict';
    const params = new URLSearchParams(window.location.search);
    const projectId = params.get("project-id");

    $(function () {
        var filter = 'all';
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');

        $('.todo-list-input').on("keydown", function (event) {
            if (event.keyCode === 13) { // Enter key
                event.preventDefault();
                $('.todo-list-add-btn').trigger("click");
            }
        });

        // Todoリストに項目を追加するための関数
        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                $.post("/todos", { name: item, projectId: projectId }, addItem)
                todoListInput.val("");
            }
        });

        // Todoリストにアイテムを追加する処理を行うaddItem()関数
        var addItem = function (item) {
            var completedClass = item.completed ? "completed" : "";
            var picture = item.picture ? item.picture : "";
            var createdAt = new Date(item.created_at);
            createdAt = createdAt.getFullYear() + '-' +
                ('0' + (createdAt.getMonth() + 1)).slice(-2) + '-' +
                ('0' + createdAt.getDate()).slice(-2) + ' ' +
                ('0' + createdAt.getHours()).slice(-2) + ':' +
                ('0' + createdAt.getMinutes()).slice(-2);
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

        // プロジェクトに紐づくTodoアイテムを取得し、それぞれのアイテムに対してaddItem()関数を呼び出して、Todoリストにアイテムを追加する処理
        $.get('/todos?projectId=' + projectId, function (items) {
            items.forEach(e => {
                addItem(e)
            });
        });

        // Todoリストのアイテムの完了を切り替えるためのjQuery関数
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

        // Todoリストのアイテムの削除を行うためのjQuery関数
        todoListItem.on('click', '.remove', function () {
            var id = $(this).closest("li").attr('id');
            var $self = $(this);
            if (confirm("このタスクを削除してもよろしいですか？")) {
                $.ajax({
                    url: "/todos/" + id,
                    type: "DELETE",
                    success: function (data) {
                        if (data.success) {
                            $self.parent().remove();
                            updateProgressBar();
                        }
                    }
                });
            }
        });

        // 完了したTodoをすべて削除するためのjQuery関数
        $('.completed-clear-btn').click(function () {
            if (confirm("完了したすべてのタスクを削除してもよろしいですか？")) {
                $.ajax({
                    url: "/todos/completed?projectId=" + projectId,
                    type: "DELETE",
                    success: function (data) {
                        if (data.success) {
                            updateProgressBar();
                            location.reload();
                        }
                    }
                });
            }
        });

        // TodoリストをフィルターするためのjQuery関数
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

        const upArrowBtn = document.querySelector('.up-arrow');
        const downArrowBtn = document.querySelector('.down-arrow');

        // Todoリストを昇順または降順にソートするためのJavaScript関数
        upArrowBtn.addEventListener('click', () => {
            upArrowBtn.style.display = 'none';
            downArrowBtn.style.display = 'block';
            var sort = 'desc';
            $.get('/todos/sorted', { filter: filter, sort: sort, projectId: projectId }, function (items) {
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
            $.get('/todos/sorted', { filter: filter, sort: sort, projectId: projectId }, function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        });

        // TodoリストをソートするためのjQuery関数
        function sortByUser() {
            $.get('/todos/sorted-by-user?projectId=' + projectId, function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        }

        function sortByCompleted() {
            $.get('/todos/sorted-by-completed?projectId=' + projectId, function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        }

        function showAll() {
            $.get('/todos?projectId=' + projectId, function (items) {
                clearList();
                items.forEach(e => {
                    addItem(e)
                });
            });
        }

        function clearList() {
            $('.todo-list').empty();
        }

        // 進捗バーを更新するためのjQuery関数
        function updateProgressBar() {
            $.get('/todos/progress?projectId=' + projectId, function (progress) {
                $('.progress-bar').css('width', progress + '%');
                $('.progress-bar').text(progress + '%');
            });
        }
    });
})(jQuery);