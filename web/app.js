const API_BASE = '/api';
const DIFFICULTIES = ['basic', 'advanced', 'expert', 'master', 're_master'];
const LEVELS = [];
for (let i = 0; i <= 15; i++) {
  LEVELS.push(String(i));
  LEVELS.push(i + '+');
}

let state = {
  fromList: [],
  genreList: [],
  levelList: LEVELS,
  filter1: { from_list: [], genre_list: [], min_ds: 1.0, max_ds: 15.9, min_level: '', max_level: '' },
  filter2: { from_list: [], genre_list: [], min_ds: 1.0, max_ds: 15.9, min_level: '', max_level: '' },
  page: 1,
  pageSize: 20,
  total: 0,
  columns: { cover: true, title: true, artist: true, ds: true, level: true, bpm: true, from: true, genre: true, type: true }
};

let debounceTimer = null;

function api(path, options = {}) {
  return fetch(API_BASE + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  }).then(r => r.json());
}

function buildFilter(player) {
  const prefix = player === 1 ? '1' : '2';
  const minDs = parseFloat(document.getElementById('minDs' + prefix).value) || 1.0;
  const maxDs = parseFloat(document.getElementById('maxDs' + prefix).value) || 15.9;
  return {
    from_list: Array.from(document.querySelectorAll('#fromList' + prefix + ' .checkbox-item.checked')).map(el => el.dataset.value),
    genre_list: Array.from(document.querySelectorAll('#genreList' + prefix + ' .checkbox-item.checked')).map(el => el.dataset.value),
    min_ds: minDs,
    max_ds: maxDs,
    min_level: document.getElementById('minLevel' + prefix).value,
    max_level: document.getElementById('maxLevel' + prefix).value
  };
}

function fetchSongs() {
  const body = {
    filter1: buildFilter(1),
    filter2: buildFilter(2),
    page: state.page,
    page_size: state.pageSize
  };
  document.getElementById('songTableBody').innerHTML = '<tr class="loading-row"><td colspan="9">加载中...</td></tr>';
  api('/select-song', { method: 'POST', body: JSON.stringify(body) })
    .then(data => {
      if (data.status === 1 && data.data) {
        state.total = data.data.total;
        renderSongList(data.data.song_list);
        renderPagination();
      } else {
        const errMsg = data.msg || data.error || JSON.stringify(data);
        document.getElementById('songTableBody').innerHTML = '<tr class="empty-row"><td colspan="9">加载失败: ' + errMsg + '</td></tr>';
      }
    })
    .catch(() => {
      document.getElementById('songTableBody').innerHTML = '<tr class="empty-row"><td colspan="9">网络错误</td></tr>';
    });
}

function renderSongList(songs) {
  const tbody = document.getElementById('songTableBody');
  if (!songs || songs.length === 0) {
    tbody.innerHTML = '<tr class="empty-row"><td colspan="9">没有找到符合条件的歌曲</td></tr>';
    return;
  }
  tbody.innerHTML = songs.map(song => {
    const dsParts = DIFFICULTIES.map(d => {
      const v = song.ds && song.ds[d];
      return v !== undefined && v !== null ? `<span class="ds-cell">${v.toFixed(1)}</span>` : '';
    }).join(' ');
    const levelParts = DIFFICULTIES.map(d => {
      const v = song.level && song.level[d];
      if (!v) return '';
      const cls = 'level-badge ' + d;
      return `<span class="${cls}">${v}</span>`;
    }).join('');
    return `<tr>
      ${col('cover', `<img class="cover-cell" src="${song.cover_url || ''}" alt="${song.title}" loading="lazy">`)}
      ${col('title', `<span class="title-cell">${song.title}</span>`)}
      ${col('artist', `<span class="artist-cell">${song.artist}</span>`)}
      ${col('ds', dsParts)}
      ${col('level', levelParts)}
      ${col('bpm', `<span class="bpm-cell">${song.bpm || '-'}</span>`)}
      ${col('from', `<span class="from-cell">${song.from || '-'}</span>`)}
      ${col('genre', `<span class="genre-cell">${song.genre}</span>`)}
      ${col('type', `<span class="type-cell">${song.type || '-'}</span>`)}
    </tr>`;
  }).join('');
}

function col(name, content) {
  return state.columns[name] ? `<td data-col="${name}">${content}</td>` : '';
}

function renderPagination() {
  const totalPages = Math.ceil(state.total / state.pageSize) || 1;
  const controls = document.getElementById('pageControls');
  const info = document.getElementById('pageInfo');

  let btns = '';
  btns += `<button class="page-btn" id="prevPage" ${state.page <= 1 ? 'disabled' : ''}>上一页</button>`;
  const maxBtns = 5;
  let startPage = Math.max(1, state.page - 2);
  let endPage = Math.min(totalPages, startPage + maxBtns - 1);
  if (endPage - startPage < maxBtns - 1) startPage = Math.max(1, endPage - maxBtns + 1);
  if (startPage > 1) btns += `<button class="page-btn" data-page="1">1</button>`;
  if (startPage > 2) btns += `<span style="color:var(--text-muted)">...</span>`;
  for (let i = startPage; i <= endPage; i++) {
    btns += `<button class="page-btn ${i === state.page ? 'active' : ''}" data-page="${i}">${i}</button>`;
  }
  if (endPage < totalPages - 1) btns += `<span style="color:var(--text-muted)">...</span>`;
  if (endPage < totalPages) btns += `<button class="page-btn" data-page="${totalPages}">${totalPages}</button>`;
  btns += `<button class="page-btn" id="nextPage" ${state.page >= totalPages ? 'disabled' : ''}>下一页</button>`;
  controls.innerHTML = btns;
  info.textContent = `共 ${state.total} 首`;

  controls.querySelectorAll('.page-btn[data-page]').forEach(btn => {
    btn.addEventListener('click', () => {
      state.page = parseInt(btn.dataset.page);
      fetchSongs();
    });
  });
  document.getElementById('prevPage')?.addEventListener('click', () => { if (state.page > 1) { state.page--; fetchSongs(); } });
  document.getElementById('nextPage')?.addEventListener('click', () => { if (state.page < totalPages) { state.page++; fetchSongs(); } });
}

function renderLevelSelects() {
  ['1', '2'].forEach(p => {
    const minEl = document.getElementById('minLevel' + p);
    const maxEl = document.getElementById('maxLevel' + p);
    let opts = '<option value="">不限</option>' + LEVELS.map(l => `<option value="${l}">${l}</option>`).join('');
    minEl.innerHTML = opts;
    maxEl.innerHTML = opts;
    minEl.value = '';
    maxEl.value = '';
  });
}

function renderCheckboxList(containerId, items, selected) {
  const container = document.getElementById(containerId);
  container.innerHTML = items.map(item => {
    const checked = selected.includes(item) ? 'checked' : '';
    return `<label class="checkbox-item ${checked}" data-value="${item}">
      <input type="checkbox" ${checked ? 'checked' : ''}>${item}
    </label>`;
  }).join('');
  container.querySelectorAll('input[type="checkbox"]').forEach(input => {
    input.addEventListener('change', () => {
      const label = input.parentElement;
      label.classList.toggle('checked', input.checked);
      triggerSearch();
    });
  });
}

function triggerSearch() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    state.page = 1;
    fetchSongs();
  }, 300);
}

function setupEventListeners() {
  ['minDs1', 'maxDs1', 'minDs2', 'maxDs2'].forEach(id => {
    document.getElementById(id).addEventListener('input', triggerSearch);
  });
  ['minLevel1', 'maxLevel1', 'minLevel2', 'maxLevel2'].forEach(id => {
    document.getElementById(id).addEventListener('change', triggerSearch);
  });

  document.getElementById('pageSizeSelect').addEventListener('change', e => {
    state.pageSize = parseInt(e.target.value);
    state.page = 1;
    fetchSongs();
  });

  document.getElementById('colToggle').querySelectorAll('.col-toggle-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const col = btn.dataset.col;
      const wasActive = btn.classList.contains('active');
      if (wasActive) {
        btn.classList.remove('active');
        state.columns[col] = false;
        document.querySelectorAll(`th[data-col="${col}"], td[data-col="${col}"]`).forEach(el => {
          el.style.display = 'none';
        });
      } else {
        btn.classList.add('active');
        state.columns[col] = true;
        document.querySelectorAll(`th[data-col="${col}"], td[data-col="${col}"]`).forEach(el => {
          el.style.display = '';
        });
      }
    });
  });

  document.getElementById('tabBar').querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      const tab = btn.dataset.tab;
      document.getElementById('filter1').classList.toggle('mobile-active', tab === '1');
      document.getElementById('filter2').classList.toggle('mobile-active', tab === '2');
    });
  });

  document.getElementById('refreshBtn').addEventListener('click', () => {
    document.getElementById('refreshBtn').textContent = '同步中...';
    api('/sync-data').then(data => {
      document.getElementById('refreshBtn').textContent = '刷新数据';
      if (data.status === 1) {
        fetchMeta();
      }
    }).catch(() => {
      document.getElementById('refreshBtn').textContent = '刷新数据';
    });
  });
}

function fetchMeta() {
  Promise.all([api('/version'), api('/from'), api('/genre'), api('/level')]).then(([versionData, fromData, genreData, levelData]) => {
    if (versionData.data?.version) {
      document.getElementById('versionDisplay').textContent = versionData.data.version;
    }
    state.fromList = fromData.data?.from_list || [];
    state.genreList = genreData.data?.genre_list || [];
    state.levelList = levelData.data?.level_list || LEVELS;
    renderCheckboxList('fromList1', state.fromList, []);
    renderCheckboxList('genreList1', state.genreList, []);
    renderCheckboxList('fromList2', state.fromList, []);
    renderCheckboxList('genreList2', state.genreList, []);
    renderLevelSelects();
    fetchSongs();
  }).catch(() => {
    document.getElementById('songTableBody').innerHTML = '<tr class="empty-row"><td colspan="9">加载元数据失败，请刷新页面</td></tr>';
  });
}

function init() {
  renderLevelSelects();
  setupEventListeners();
  fetchMeta();
}

init();
