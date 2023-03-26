(function ($) {
    'use strict';

    var originalItems;
    const nameHeader = document.getElementById("sort-name");
    const codeHeader = document.getElementById("sort-code");
    const priorityHeader = document.getElementById("sort-priority");
    const colorHeader = document.getElementById("sort-color");
    const initialClass = "d-none";
    nameHeader.querySelector(".fa-chevron-up").classList.add(initialClass);
    codeHeader.querySelector(".fa-chevron-up").classList.add(initialClass);
    priorityHeader.querySelector(".fa-chevron-up").classList.add(initialClass);
    colorHeader.querySelector(".fa-chevron-up").classList.add(initialClass);
    var availableUsersLoaded = false;
    var participantsLoaded = false;

    $(function () {
        var projectListItem = $('.project-list');
        var projectName = $('#project-name');
        var projectCode = $('#project-code');
        var projectDescription = $('#project-description');
        var projectColor = $('#project-color');
        var projectPriority = $('#project-priority');
        var sortOrder = 'asc';

        $('#project-create-btn').on("click", function () {
            var name = projectName.val();
            var code = projectCode.val();
            var description = projectDescription.val();
            var selectedColor = $('.color-option.active').data("color");
            var priority = projectPriority.val();

            if (!name || !code) {
                alert("プロジェクト名とコードを入力してください。");
                return;
            }

            // 新しいプロジェクトを作成するためのjQuery関数
            $.post('/projects', {
                name: name,
                code: code,
                description: description,
                color: selectedColor,
                priority: priority
            }).done(function () {
                projectName.val("");
                projectCode.val("");
                projectDescription.val("");

                location.reload();
            });
        });

        $.ajaxSetup({
            cache: true
        });

        // プロジェクトのリストを取得し、各プロジェクトのユーザー情報を取得してリストアイテムを生成するためのjQuery関数
        var createListItemHtml = function (item) {
            var color = item.color || "#A9A9A9";
            var priorityText = "";
            var listColor = "";
            switch (item.priority) {
                case "high":
                    priorityText = "上";
                    listColor = "pink";
                    break;
                case "mid":
                    priorityText = "中";
                    break;
                case "low":
                    priorityText = "下";
                    break;
                default:
                    priorityText = "";
                    break;
            }
            var listColorStyle = listColor ? "style='background-color: " + listColor + ";'" : "";

            return "<li class='project-item'" + listColorStyle + " data-id='" + item.id + "'>" +
                "<div class='project-color' style='background-color: " + color + ";'></div>" +
                "<div class='project-name'>" + item.name + "</div>" +
                "<div class='project-code'>&nbsp;(" + item.code + ")</div>" +
                "<div class='project-priority'>" + priorityText + "</div>" +
                "<img class='project-picture' src='" + item.user_picture + "'/>" +
                "</li>";
        };

        $.get('/projects', function (items) {
            originalItems = items;
            var users = [];
            items.forEach(function (item) {
                users.push($.get("/user/" + item.user_id));
            });
            if (items.length === 1) {
                $.when.apply(null, users).done(function (user) {
                    items[0].user_picture = user.picture || "";
                    addItem(items[0]);
                });
            } else {
                $.when.apply(null, users).done(function () {
                    for (var i = 0; i < arguments.length; i++) {
                        var user = arguments[i][0];
                        items[i].user_picture = user.picture || "";
                        addItem(items[i]);
                    }
                });
            }
        });

        var addItem = function (item) {
            var listItemHtml = createListItemHtml(item);
            projectListItem.append(listItemHtml);
        };

        $("#color-options .color-option").click(function () {
            $(".color-option").removeClass("active");
            $(this).addClass("active");
            projectColor.val($(this).data("color"));
        });

        $("#color-options-modal .color-option").click(function () {
            $(".color-option").removeClass("active");
            $(this).addClass("active");
            projectColor.val($(this).data("color"));
        });

        $('#search-project').keyup(function () {
            var searchValue = $(this).val().toLowerCase();
            if (searchValue === '') {
                $('.project-item').show();
            } else {
                searchProjects(searchValue);
            }
        });

        nameHeader.addEventListener("click", function () {
            toggleChevronClass(this);
            sortItems(originalItems, 'name');
        });

        codeHeader.addEventListener("click", function () {
            toggleChevronClass(this);
            sortItems(originalItems, 'code');
        });

        priorityHeader.addEventListener("click", function () {
            toggleChevronClass(this);
            sortItems(originalItems, 'priority');
        });

        colorHeader.addEventListener("click", function () {
            toggleChevronClass(this);
            sortItems(originalItems, 'color');
        });

        // 配列のオブジェクトを指定したキーでソートするためのjQuery関数
        function sortItems(items, sortBy) {
            items.sort(function (a, b) {
                var aValue = a[sortBy].toLowerCase();
                var bValue = b[sortBy].toLowerCase();
                if (sortOrder === 'asc') {
                    return aValue.localeCompare(bValue);
                } else {
                    return bValue.localeCompare(aValue);
                }
            });
            $('.project-list').empty();
            items.forEach(function (item) {
                addItem(item);
            });
            sortOrder = (sortOrder === 'asc') ? 'desc' : 'asc';
        }

        $(".project-list").on("click", "li.project-item", function () {
            var itemId = $(this).data("id");
            $("#edit-project-form").hide();
            $("#show-edit-project-btn").hide();

            // 特定のプロジェクトの情報を取得するためのjQuery関数
            $.get(`/projects/${itemId}`)
                .done(function (project) {
                    var modal = $("#project-detail-modal");
                    var projectTitle = project.name + " (" + project.code + ")";
                    var participantsList = [];
                    var availableUsersList = [];

                    // プロジェクトを編集できるかどうかを確認するためのjQuery関数
                    $.get(`/projects/${itemId}/check-edit-auth`, function (response) {
                        if (response) {
                            $("#show-edit-project-btn").show().on("click", function () {
                                $("#edit-project-form").show();
                                $("#project-name-modal").val(project.name);
                                $("#project-code-modal").val(project.code);
                                $("#project-priority-modal").val(project.priority);
                                var color = project.color || "#A9A9A9";
                                $("#color-options-modal .color-option").removeClass("active");
                                $("#color-options-modal .color-option[data-color='" + color + "']").addClass('active');
                                $("#project-description-text").val(project.description);
                                if (!participantsLoaded) {
                                    var select = $('#project-manager');
                                    $.each(participantsList, function (_, participant) {
                                        select.append($('<option>', {
                                            value: participant['id'],
                                            text: participant['email']
                                        }));
                                    });
                                    select.val(project.user_id);
                                    participantsLoaded = true
                                }
                                var availableUsersSelect = $('#available-users');
                                var selectedUsersSelect = $('#selected-users');
                                if (!availableUsersLoaded) {

                                    // プロジェクトの利用可能なユーザーを取得し、フォーム内のユーザー選択リストを更新するためのjQuery関数
                                    $.get(`/projects/${itemId}/availableUsers`, function (availableUsers) {
                                        availableUsers.forEach(function (availableUser) {
                                            if (availableUser.id != project.user_i) {
                                                if (!availableUsersList.some(user => user.id === availableUser.id)) {
                                                    availableUsersList.push({ id: availableUser.id, email: availableUser.email });
                                                    availableUsersSelect.append($('<option>', {
                                                        value: availableUser.id,
                                                        text: availableUser.email
                                                    }));
                                                }
                                            }
                                        });
                                        participantsList.forEach(function (participant) {
                                            if (participant.id != project.user_id) {
                                                selectedUsersSelect.append($('<option>', {
                                                    value: participant.id,
                                                    text: participant.email
                                                }));
                                            }
                                        });
                                    })
                                        .fail(function (errorThrown) {
                                            console.error(`Failed to load available users: ${errorThrown}`);
                                        })
                                        .always(function () {
                                            availableUsersLoaded = true;
                                        });
                                }
                                $('#add-user-btn').click(function () {
                                    var selectedOptions = availableUsersSelect.find(':selected');
                                    selectedOptions.clone().appendTo(selectedUsersSelect);
                                    selectedOptions.remove();
                                });
                                $('#remove-user-btn').click(function () {
                                    var selectedOptions = selectedUsersSelect.find(':selected');
                                    selectedOptions.clone().appendTo(availableUsersSelect);
                                    selectedOptions.remove();
                                });
                            });

                            // プロジェクトを削除するためのjQuery関数
                            $("#delete-project-link").on("click", function () {
                                if (confirm("このプロジェクトを削除しますか？")) {
                                    $.ajax({
                                        url: `/projects/${itemId}`,
                                        type: "DELETE",
                                        success: function () {
                                            alert("プロジェクトが削除されました。");
                                            location.reload();
                                        },
                                        error: function () {
                                            alert("プロジェクトの削除に失敗しました。");
                                        }
                                    });
                                }
                            });

                            // プロジェクト情報を更新するためのjQuery関数
                            $("#edit-project-btn").on("click", function () {
                                const confirmPromise = new Promise((resolve, reject) => {
                                    if (confirm("プロジェクト情報を更新しますか？")) {
                                        resolve();
                                    } else {
                                        reject();
                                    }
                                });
                                confirmPromise.then(() => {
                                    $.ajax({
                                        url: `/projects/${itemId}`,
                                        type: "PUT",
                                        data: {
                                            name: $("#project-name-modal").val(),
                                            code: $("#project-code-modal").val(),
                                            description: $("#project-description-text").val(),
                                            color: $("#color-options-modal .color-option.active").data("color"),
                                            priority: $("#project-priority-modal").val(),
                                            managerId: $("#project-manager").val(),
                                            participantIds: $("#selected-users option").map(function () { return $(this).val(); }).get(),
                                            availableUserIds: $("#available-users option").map(function () { return $(this).val(); }).get()
                                        },
                                        success: function () {
                                            alert("プロジェクトが更新されました。");
                                            location.reload();
                                        },
                                        error: function () {
                                            alert("プロジェクトの更新に失敗しました。");
                                            location.reload();
                                        }
                                    });
                                }).catch(() => {
                                });
                            });
                            $("#modal-close-btn").on("click", function () {
                                $("#edit-project-form").hide();
                            });
                        } else {
                            $("#show-project-btn").hide();
                        }
                    }).fail(function () {
                        alert("プロジェクトの編集権限チェックに失敗しました。");
                        $("#edit-project-btn").hide();
                    });

                    modal.find(".modal-title").text(projectTitle);
                    var headerColor = project.color || "#A9A9A9"; // 기본값으로 회색 지정
                    modal.find(".modal-header").attr("style", "background-color: " + headerColor + ";");
                    modal.find(".modal-footer").attr("style", "background-color: " + headerColor + ";");
                    var priorityText = "";
                    switch (project.priority) {
                        case "high":
                            priorityText = "上";
                            break;
                        case "mid":
                            priorityText = "中";
                            break;
                        case "low":
                            priorityText = "下";
                            break;
                        default:
                            priorityText = "";
                            break;
                    }
                    modal.find("#project-priority").text(priorityText)

                    // プロジェクトのユーザー情報を取得し、モーダルで表示するためのjQuery関数
                    $.get("/user/" + project.user_id, function (user) {
                        modal.find("#project-userId").text(user.email);
                        if (user.picture) {
                            modal.find("#project-userPicture").attr("src", user.picture);
                        } else {
                            modal.find("#project-userPicture").hide();
                        }
                    });

                    modal.find("#project-description").html(project.description.replace(/\n/g, "<br>"));
                    modal.find("#project-createdAt").text(project.created_at);

                    // プロジェクトの参加者情報を取得し、モーダルで表示するためのjQuery関数
                    $.get(`/projects/${itemId}/participants`, function (participants) {
                        var participantList = ""
                        participants.forEach(function (participant) {
                            participantsList.push({ id: participant.id, email: participant.email });
                            if (participant.id != project.user_id) {
                                participantList += `
                            <li>
                                <div class="participant">
                                    <img src="${participant.picture}" class="rounded-circle">
                                    <div>${participant.email}</div>
                                </div>
                            </li>`;
                            }
                        });
                        $(".participant-list ul").html(participantList);
                    });
                    $("#join-project-btn").on("click", function () {
                        window.location.href = "/projects/" + itemId + "/todos";
                    });
                    modal.modal("show");
                })
                .fail(function () {
                    alert("プロジェクト情報の取得に失敗しました。");
                    $("#project-detail-modal").modal("hide");
                });
        });
    });

    // プロジェクトアイテムを取得し、検索語が含まれている場合は表示し、含まれていない場合は非表示にするjQuery関数
    function searchProjects(searchValue) {
        $('.project-item').each(function () {
            var itemName = $(this).find('.project-name').text().toLowerCase();
            var itemCode = $(this).find('.project-code').text().toLowerCase();
            if (itemName.indexOf(searchValue) >= 0 || itemCode.indexOf(searchValue) >= 0) {
                $(this).show();
            } else {
                $(this).hide();
            }
        });
    }

    // 要素の子要素のFontAwesomeアイコン要素をトグルするjQuery関数
    function toggleChevronClass(element) {
        const up = "fa-chevron-up";
        const down = "fa-chevron-down";

        const chevronUp = element.querySelector(`.${up}`);
        const chevronDown = element.querySelector(`.${down}`);

        if (chevronUp.classList.contains("d-none")) {
            chevronUp.classList.remove("d-none");
            chevronDown.classList.add("d-none");
        } else {
            chevronUp.classList.add("d-none");
            chevronDown.classList.remove("d-none");
        }
    }
})(jQuery);