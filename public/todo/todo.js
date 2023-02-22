(function ($) {
    'use strict';
    $(function () {
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
            console.log(item.picture);
            var completedClass = item.completed ? "completed" : "";
            var picture = item.picture ? item.picture : "";
            var createdAt = item.created_at;
            var listItemHtml =
                "<li class='" +
                completedClass +
                "' id='" +
                item.id +
                "' style='display: flex; justify-content: space-between;'>" +
                "<div class='profile-image-container' style='margin-right: 10px;'>" +
                "<img class='profile-image' src='" + picture + "' style='width:30px; height:30px;'/>" +
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
            })
        });

        todoListItem.on('click', '.remove', function () {
            var id = $(this).closest("li").attr('id')
            var $self = $(this);
            $.ajax({
                url: "/todos/" + id,
                type: "DELETE",
                success: function (data) {
                    if (data.success) {
                        $self.parent().remove();
                    }
                }
            })
        });
    });
})(jQuery);