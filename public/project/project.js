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

        var addItem = function (item) {
            console.log(item.user_id)
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
                    "<li class='project-item'" + listColorStyle + ">" +
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

    function updateProjectList(projects) {
        $('.project-list').html('');
        projects.forEach(e => {
            addItem(e);
        });
    }
})(jQuery);