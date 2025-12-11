// ---------- 全局状态 ----------
const BASE_URL = "http://localhost:8080";
const DELETE_LOG_URL = BASE_URL + "/v1/log";
const DELETE_LOG_BY_TASK_ID_URL = BASE_URL + "/v1/logs/task";
const GET_LOG_URL = BASE_URL + "/v1/log";
const LIST_LOG_URL = BASE_URL + "/v1/logs";
const LIST_LOG_BY_TASK_ID_URL = BASE_URL + "/v1/logs/task";
const CREATE_ONCE_URL = BASE_URL + "/v1/once";
const GET_ONCE_LOG_URL = BASE_URL + "/v1/once/log";
const DELETE_ONCE_URL = BASE_URL + "/v1/once";
const LIST_ONCE_URL = BASE_URL + "/v1/onces";
const CREATE_TASK_URL = BASE_URL + "/v1/task";
const DELETE_TASK_URL = BASE_URL + "/v1/task";
const UPDATE_TASK_URL = BASE_URL + "/v1/task";
const GET_TASK_URL = BASE_URL + "/v1/task";
const LIST_TASK_URL = BASE_URL + "/v1/tasks";

let currentPage = 0;
const limit = 10;

// ---------- Tab 切换 ----------
document.querySelectorAll(".tab-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    document
      .querySelectorAll(".tab-content")
      .forEach((el) => el.classList.add("hidden"));
    document
      .querySelectorAll(".tab-btn")
      .forEach((b) => b.classList.remove("font-bold", "text-blue-600"));
    const tab = btn.dataset.tab;
    document.getElementById(`tab-${tab}`).classList.remove("hidden");
    btn.classList.add("font-bold", "text-blue-600");
  });
});

// ---------- 一次性任务 ----------
document.getElementById("form-once").addEventListener("submit", async (e) => {
  e.preventDefault();
  const source = document.getElementById("once-source").value;
  const dest = document.getElementById("once-dest").value;
  const res = await fetch(CREATE_ONCE_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ source, destination: dest }),
  });
  const data = await res.json();
  if (data.success) {
    if (data.data.id) {
      const logContainer = document.getElementById("once-log-container");
      logContainer.innerHTML = "";
      onceEventSource = new EventSource(
        `${GET_ONCE_LOG_URL}?id=${data.data.id}`
      );
      onceEventSource.onmessage = (ev) => {
        const log = ev.data;
        const safeLog = log
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")
          .replace(/\n/g, "<br>");
        logContainer.innerHTML += safeLog + "<br>";
        logContainer.scrollTop = logContainer.scrollHeight;
      };
    }
    loadOnceTaskIds(); // 自动刷新下拉
  } else {
    alert("创建失败: " + data.error);
  }
});

async function loadOnceTaskIds() {
  const res = await fetch(LIST_ONCE_URL);
  const onceList = await res.json();
  const { data } = onceList;
  const select = document.getElementById("once-id-select");
  select.innerHTML = '<option value="">请选择任务 ID</option>';
  data.forEach((item) => {
    const opt = document.createElement("option");
    opt.value = item.id;
    opt.textContent = `#${item.id} ${item.source} → ${item.dest}`;
    select.appendChild(opt);
  });
}

document
  .getElementById("btn-load-once-ids")
  .addEventListener("click", loadOnceTaskIds);

let onceEventSource = null;
// 一次性任务日志 SSE
document.getElementById("once-id-select").addEventListener("change", (e) => {
  const id = e.target.value;
  if (onceEventSource) onceEventSource.close();
  if (!id) return;
  const logContainer = document.getElementById("once-log-container");
  logContainer.innerHTML = "";
  onceEventSource = new EventSource(`${GET_ONCE_LOG_URL}?id=${id}`);
  onceEventSource.onmessage = (ev) => {
    const log = ev.data;
    const safeLog = log
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/\n/g, "<br>");
    logContainer.innerHTML += safeLog + "<br>";
    logContainer.scrollTop = logContainer.scrollHeight;
  };
});
loadOnceTaskIds(); // 初始加载

// ---------- 定时任务 ----------
document.getElementById("form-task").addEventListener("submit", async (e) => {
  e.preventDefault();
  const cron = document.getElementById("task-cron").value;
  const source = document.getElementById("task-source").value;
  const dest = document.getElementById("task-dest").value;
  const res = await fetch(CREATE_TASK_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ cron, source, destination: dest }),
  });
  const data = await res.json();
  if (data.success) {
    alert("定时任务创建成功");
    loadTasks();
  } else {
    alert("失败: " + data.error);
  }
});

async function loadTasks() {
  const res = await fetch(
    `${LIST_TASK_URL}?skip=${currentPage * limit}&limit=${limit}`
  );
  const data = await res.json();
  if (!data.success) return;
  const tbody = document.getElementById("task-list");
  tbody.innerHTML = "";
  data.data.items.forEach((t) => {
    const tr = document.createElement("tr");
    tr.className = "border-b";
    tr.innerHTML = `
    <td>${t.id}</td>
    <td>${t.cron}</td>
    <td>${t.source}</td>
    <td>${t.destination}</td>
    <td class="space-x-2">
      <button class="text-blue-600 text-sm update-btn" data-id="${t.id}" data-cron="${t.cron}" data-source="${t.source}" data-dest="${t.destination}">更新</button>
      <button class="text-red-600 delete-btn" data-id="${t.id}">删除</button>
    </td>`;
    tbody.appendChild(tr);
  });
  document.getElementById("page-info").textContent = `第 ${currentPage + 1} 页`;
  // 绑定更新按钮
  document.querySelectorAll(".update-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const id = btn.dataset.id;
      document.getElementById("update-task-id").value = id;
      document.getElementById("update-task-cron").value = btn.dataset.cron;
      document.getElementById("update-task-source").value = btn.dataset.source;
      document.getElementById("update-task-dest").value = btn.dataset.dest;
      document.getElementById("update-modal").classList.remove("hidden");
    });
  });
  // 绑定删除
  document.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      if (!confirm("确认删除？")) return;
      const id = btn.dataset.id;
      await fetch(`${DELETE_TASK_URL}?id=${id}`, { method: "DELETE" });
      loadTasks();
    });
  });
}

document.getElementById("form-update-task").addEventListener("submit", async (e) => {
    e.preventDefault();
    const id = parseInt(document.getElementById("update-task-id").value);
    const cron = document.getElementById("update-task-cron").value.trim();
    const source = document.getElementById("update-task-source").value.trim();
    const dest = document.getElementById("update-task-dest").value.trim();

    if (!cron || !source || !dest) {
      alert("请填写完整信息");
      return;
    }

    const task = { id, cron, source, destination: dest };
    const res = await fetch(`${UPDATE_TASK_URL}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(task),
    });
    const data = await res.json();
    if (data.success) {
      alert("更新成功");
      document.getElementById("update-modal").classList.add("hidden");
      loadTasks();
    } else {
      alert("更新失败: " + data.error);
    }
  });

// 取消按钮
document.getElementById("btn-cancel-update").addEventListener("click", () => {
  document.getElementById("update-modal").classList.add("hidden");
});

document.getElementById("btn-prev").addEventListener("click", () => {
  if (currentPage > 0) {
    currentPage--;
    loadTasks();
  }
});
document.getElementById("btn-next").addEventListener("click", () => {
  currentPage++;
  loadTasks();
});
document
  .getElementById("btn-refresh-tasks")
  .addEventListener("click", loadTasks);
loadTasks(); // 初始加载

// ---------- 定时任务日志 ----------
async function loadTaskIdsForLog() {
  const res = await fetch(`${LIST_TASK_URL}?limit=1000`); // 获取全部用于下拉
  const data = await res.json();
  const select = document.getElementById("task-id-select");
  select.innerHTML = '<option value="">请选择任务 ID</option>';
  data.data.items.forEach((t) => {
    const opt = document.createElement("option");
    opt.value = t.id;
    opt.textContent = `#${t.id} ${t.source} → ${t.destination}`;
    select.appendChild(opt);
  });
}

document
  .getElementById("btn-load-task-ids")
  .addEventListener("click", loadTaskIdsForLog);
loadTaskIdsForLog();

document.getElementById("task-id-select").addEventListener("change", (e) => {
  const id = e.target.value;
  if (!id) return;
  const logContainer = document.getElementById("task-log-container");
  logContainer.innerHTML = "";
  fetch(`${GET_LOG_URL}?id=${id}`).then(async (res) => {
    const data = await res.json();
    if (!data.success) return;
    const { msg: log } = data.data;
    const safeLog = log
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/\n/g, "<br>");
    logContainer.innerHTML = safeLog;
  });
});
