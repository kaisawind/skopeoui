// Tab切换
function showTab(tabName) {
  document
    .querySelectorAll(".tab-content")
    .forEach((el) => el.classList.add("hidden"));
  document
    .querySelectorAll(".tab-btn")
    .forEach((btn) => btn.classList.remove("active"));
  document.getElementById(`tab-${tabName}`).classList.remove("hidden");
  document.querySelector(`[data-tab="${tabName}"]`).classList.add("active");

  if (tabName === "task") listTasks();
  if (tabName === "log") listAllLogs();
}

// —————— Once 相关 ——————
let currentLogReader = null;
let currentLogId = null;

// 加载一次性任务列表
async function loadOnceList() {
  const res = await fetch("/v1/onces");
  const tasks = await res.json();
  const select = document.getElementById("once-select");
  select.innerHTML = '<option value="">-- 请选择一个任务 --</option>';
  tasks.forEach((task) => {
    const opt = document.createElement("option");
    opt.value = task.id;
    opt.textContent = `${task.id} | ${task.source} → ${task.dest}`;
    select.appendChild(opt);
  });
}

function refreshOnceList() {
  loadOnceList();
}

// 下拉选择变更
function onOnceSelect() {
  const id = document.getElementById("once-select").value;
  document.getElementById("btn-start-log").disabled = !id;
}

// 启动日志
function startOnceLog() {
  const id = document.getElementById("once-select").value;
  if (!id) return;
  viewOnceLog(id);
}

// 查看 SSE 日志（使用 fetch + ReadableStream，因接口是 POST）
async function viewOnceLog(id) {
  stopOnceLog();
  const logEl = document.getElementById("once-log-output");
  logEl.innerHTML = `⏳ 正在连接日志流 (ID: ${id})...\n`;
  logEl.scrollTop = logEl.scrollHeight;

  try {
    const res = await fetch(`/v1/once/log?id=${encodeURIComponent(id)}`, {
      method: "POST",
    });
    if (!res.ok) {
      logEl.innerHTML += `❌ 日志请求失败: ${res.status} ${res.statusText}\n`;
      return;
    }

    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    currentLogReader = reader;
    currentLogId = id;

    logEl.innerHTML += `✅ 连接成功，开始接收日志...\n`;
    logEl.scrollTop = logEl.scrollHeight;

    const read = async () => {
      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;
          const text = decoder.decode(value, { stream: true });
          logEl.innerHTML += text;
          logEl.scrollTop = logEl.scrollHeight;
        }
        logEl.innerHTML += `\nℹ️ 日志流已结束\n`;
      } catch (err) {
        if (err.name !== "AbortError") {
          logEl.innerHTML += `\n📡 读取错误: ${err.message}\n`;
        }
      }
    };
    read();
  } catch (err) {
    logEl.innerHTML += `\n💥 启动失败: ${err.message}\n`;
  }
}

// 停止日志
function stopOnceLog() {
  if (currentLogReader) {
    currentLogReader.cancel();
    currentLogReader = null;
    currentLogId = null;
  }
  const logEl = document.getElementById("once-log-output");
  if (logEl.innerHTML.trim() === "") {
    logEl.innerHTML = "日志已停止。";
  }
}

// 创建一次性任务（可选）
async function createOnce() {
  const source = document.getElementById("once-source").value;
  const dest = document.getElementById("once-dest").value;
  if (!source || !dest) {
    alert("请输入源和目标镜像");
    return;
  }
  const res = await fetch("/v1/once", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ source, destination: dest }),
  });
  if (res.ok) {
    alert("任务已提交");
    document.getElementById("once-source").value = "";
    document.getElementById("once-dest").value = "";
    loadOnceList(); // 自动刷新列表
  } else {
    alert("提交失败: " + (await res.text()));
  }
}

// 初始化
loadOnceList();

// —————— Task 相关 ——————
async function createTask() {
  const cron = document.getElementById("task-cron").value;
  const source = document.getElementById("task-source").value;
  const dest = document.getElementById("task-dest").value;
  await fetch("/v1/task", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ cron, source, destination: dest }),
  });
  listTasks();
}

async function updateTask() {
  const cron = document.getElementById("task-cron").value;
  const source = document.getElementById("task-source").value;
  const dest = document.getElementById("task-dest").value;
  // 假设 id 在 source 字段中不现实，此处简化：前端需先查任务再更新
  alert("更新需先获取任务详情，此处略");
}

async function listTasks() {
  const skip = document.getElementById("task-skip").value || 0;
  const limit = document.getElementById("task-limit").value || 10;
  const res = await fetch(`/v1/tasks?skip=${skip}&limit=${limit}`);
  const { items } = await res.json();
  const listEl = document.getElementById("task-list");
  listEl.innerHTML = items
    .map(
      (t) => `
    <div class="flex justify-between items-center bg-gray-100 p-3 rounded">
      <div>ID: ${t.id} | ${t.cron} | ${t.source} → ${t.destination}</div>
      <button onclick="deleteTask(${t.id})" class="text-red-600">删除</button>
    </div>
  `
    )
    .join("");
}

async function deleteTask(id) {
  if (!confirm("确认删除？")) return;
  await fetch(`/v1/task?id=${id}`, { method: "DELETE" });
  listTasks();
}

// —————— Log 相关 ——————
async function listAllLogs() {
  const res = await fetch("/v1/logs?limit=50");
  const { items } = await res.json();
  renderLogs(items);
}

async function listLogsByTask() {
  const taskId = document.getElementById("log-taskId").value;
  if (!taskId) return alert("请输入任务ID");
  const res = await fetch(`/v1/logs/task?taskId=${taskId}&limit=50`);
  const { items } = await res.json();
  renderLogs(items);
}

function renderLogs(logs) {
  const logEl = document.getElementById("log-list");
  logEl.innerHTML = logs
    .map(
      (log) =>
        `<div class="log-line">[${new Date(log.time / 1e6).toISOString()}] ${
          log.msg
        }</div>`
    )
    .join("");
}
