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

        $(".project-list").on("click", "li.project-item", function () {
            var itemId = $(this).data("id");
            $.get(`/projects/${itemId}`, function (project) {
                var modal = $("#project-detail-modal");
                modal.find(".modal-title").text(project.name);
                modal.find("#project-id").text(project.id);
                modal.find("#project-name").text(project.name);
                modal.find("#project-code").text(project.code);
                modal.find("#project-color").text(project.color);
                modal.find("#project-description").text(project.description);
                modal.find("#project-priority").text(project.priority);
                modal.find("#project-userId").text(project.user_id);
                modal.find("#project-createdAt").text(project.created_at);

                $.get(`/projects/${itemId}/participants`, function (participants) {
                    // 참여자 목록을 리스트 형태로 만들어서 모달창에 채우기
                    var participantList = "<ul>";
                    participants.forEach(function (participant) {
                        participantList += `<li>${participant.name} (${participant.email})</li>`;
                    });
                    participantList += "</ul>";
                    $("#project-participants").html(participantList);
                });
            });
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
                //     // 삭제 버튼 클릭 시 프로젝트를 삭제함
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


            $("#project-detail-modal").modal("show");
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