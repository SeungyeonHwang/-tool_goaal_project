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
                alert("Please fill in project name and code.");
                return;
            }
            $.post('/projects', {
                name: name,
                code: code,
                description: description,
                color: selectedColor,
                priority: priority
            });
            projectName.val("");
            projectCode.val("");
            projectDescription.val("");
        });

        $.ajaxSetup({
            cache: true
        });

        var addItem = function (item) {
            var color = item.color || "#A9A9A9";
            var priorityText = ""; // 우선순위에 따른 한자 표시를 위한 문자열 변수
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

            $.get("/user/" + item.user_id, function (user) {
                var pictureUrl = user.picture || "";

                var listItemHtml =
                    "<li class='project-item'" + listColorStyle + " data-id='" + item.id + "'>" +
                    "<div class='project-color' style='background-color: " + color + ";'></div>" +
                    "<div class='project-name'>" + item.name + "</div>" +
                    "<div class='project-code'>&nbsp;(" + item.code + ")</div>" +
                    "<div class='project-priority'>" + priorityText + "</div>" +
                    "<img class='project-picture' src='" + pictureUrl + "'/>" +
                    "</li>";
                projectListItem.append(listItemHtml);
            });
        };

        $.get('/projects', function (items) {
            originalItems = items
            items.forEach(e => {
                addItem(e);
            });
        });

        $("#color-options .color-option").click(function () {
            $(".color-option").removeClass("active");
            $(this).addClass("active");
            projectColor.val($(this).data("color"));
        });

        $('#search-project').keyup(function () {
            var searchValue = $(this).val().toLowerCase();
            if (searchValue === '') {
                $('.project-item').show();
                $('.project-list').html('');
                originalItems.forEach(e => {
                    addItem(e);
                });
            } else {
                searchProjects(searchValue);
            }
        });

        const nameHeader = document.getElementById("sort-name");
        const codeHeader = document.getElementById("sort-code");
        const priorityHeader = document.getElementById("sort-priority");
        const colorHeader = document.getElementById("sort-color");

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

        // 프로젝트 상세 모달 출력
        $(".project-list").on("click", "li.project-item", function () {
            var itemId = $(this).data("id");

            // AJAX 요청
            $.get(`/projects/${itemId}`)
                .done(function (project) {
                    var modal = $("#project-detail-modal");
                    var projectTitle = project.name + " (" + project.code + ")";
                    var participantsList = [];
                    var availableUsersList = [];                

                    // Edit 버튼 표시 여부 설정
                    $.get(`/projects/${itemId}/check-edit-auth`, function (response) {
                        if (response) {
                            $("#edit-project-form").hide();
                            $("#show-edit-project-btn").show().on("click", function () {
                                // 각 폼 필드에 해당하는 프로젝트 상세 정보 적용하기
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
                                    $.get(`/projects/${itemId}/availableUsers`, function (availableUsers) {
                                        availableUsers.forEach(function (availableUser) {
                                            availableUsersList.push({ id: availableUser.id, email: availableUser.email });
                                            availableUsersSelect.append($('<option>', {
                                                value: availableUser.id,
                                                text: availableUser.email
                                            }));
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
                                        .always(function() {
                                            // $.get() 함수가 완료되면 플래그 변수를 true로 설정합니다.
                                            availableUsersLoaded = true;
                                        });
                                }
                                // add-user-btn 버튼 클릭 이벤트 처리
                                $('#add-user-btn').click(function () {
                                    var selectedOptions = availableUsersSelect.find(':selected');
                                    selectedOptions.clone().appendTo(selectedUsersSelect);
                                    selectedOptions.remove();
                                });
                                // remove-user-btn 버튼 클릭 이벤트 처리
                                $('#remove-user-btn').click(function () {
                                    var selectedOptions = selectedUsersSelect.find(':selected');
                                    selectedOptions.clone().appendTo(availableUsersSelect);
                                    selectedOptions.remove();
                                });
                            });
                            $("#edit-project-btn").on("click", function () {
                                // HTTP PUT 요청
                                $.ajax({
                                    url: `/projects/${itemId}`,
                                    type: "PUT",
                                    data: {
                                        // 프로젝트 정보 업데이트에 필요한 데이터
                                    },
                                    success: function () {
                                        console.log("프로젝트 정보가 업데이트되었습니다.");
                                    },
                                    error: function () {
                                        console.log("프로젝트 정보 업데이트에 실패했습니다.");
                                    }
                                });
                            });

                            // 모달 숨김 버튼 클릭 시 이벤트 처리
                            $("#modal-close-btn").on("click", function () {
                                $("#edit-project-form").hide();
                            });
                        } else {
                            $("#show-project-btn").hide();
                        }
                    }).fail(function () {
                        // 요청 실패 시 실행할 코드
                        alert("프로젝트 수정 권한 체크에 실패했습니다.");
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

                    // 참여자 목록 출력
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

                    modal.modal("show");

                })
                .fail(function () {
                    // 요청 실패 시 실행할 코드
                    alert("프로젝트 정보를 가져오는데 실패했습니다.");
                    $("#project-detail-modal").modal("hide");
                });
            /////

            //TODO : 가져오는거 성공했을 경우에만
            console.log(sessionStorage)

            // // 모달창에 수정 및 삭제 버튼에 클릭 이벤트 추가하기
            // $("#project-edit-btn").off("click").on("click", function () {
            //     // 수정 버튼 클릭 시 프로젝트 정보를 가지고 프로젝트 수정 모달창을 보여줌
            //     $("#project-edit-modal #project-id").val(project.id);
            //     $("#project-edit-modal #project-name").val(project.name);
            //     $("#project-edit-modal #project-code").val(project.code);
            //     $("#project-edit-modal #project-description").val(project.description);
            //     $("#project-edit-modal .color-option").removeClass("active");
            //     $("#project-edit-modal .color-option[data-color='" + project.color + "']").addClass("active");
            //     $("#project-edit-modal #project-priority").val(project.priority);
            //     $("#project-detail-modal").modal("hide");
            //     $("#project-edit-modal").modal("show");
            // });

            // $("#project-delete-btn").off("click").on("click", function () {
                // 삭제 버튼 클릭 시 프로젝트를 삭제함
            //     if (confirm("Are you sure you want to delete this project?")) {
            //         $.ajax({
            //             url: `/projects/${project.id}`,
            //             type: 'DELETE',
            //             success: function () {
            //                 // 삭제 성공 시 프로젝트 리스트를 업데이트함
            //                 $.get('/projects', function (projects) {
            //                     updateProjectList(projects);
            //                 });
            //             }
            //         });
            //         $("#project-detail-modal").modal("hide");
            //     }
            // });

        });
    });

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