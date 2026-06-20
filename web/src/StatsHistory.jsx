import { useState, useEffect, useCallback } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { getStatsHistory, getUsers } from "./api";

const fmt = (ts) =>
  new Date(ts).toLocaleString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });

// Экран истории статов -> GET /api/v1/user/stats/history.
// Выбор пользователя — выпадающий список из GET /api/v1/users.
export default function StatsHistory({ usersVersion = 0 }) {
  const [users, setUsers] = useState([]);
  const [usersError, setUsersError] = useState(null);
  const [userId, setUserId] = useState("");
  const [limit, setLimit] = useState("50");
  const [rows, setRows] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [loaded, setLoaded] = useState(false);

  // Загружаем список пользователей при монтировании и после добавления нового (usersVersion).
  useEffect(() => {
    let cancelled = false;
    getUsers()
      .then((data) => {
        if (!cancelled) {
          setUsers(data || []);
          setUsersError(null);
        }
      })
      .catch((err) => {
        if (!cancelled) setUsersError(err.message);
      });
    return () => {
      cancelled = true;
    };
  }, [usersVersion]);

  const load = useCallback(
    async (uid, lim) => {
      if (!uid) return;
      setLoading(true);
      setError(null);
      try {
        const data = await getStatsHistory(uid, Number(lim) || 0);
        setRows(data || []);
        setLoaded(true);
      } catch (err) {
        setError(err.message);
        setRows([]);
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // Автоподгрузка истории при выборе пользователя (и при смене лимита).
  useEffect(() => {
    if (userId) {
      load(userId, limit);
    } else {
      setRows([]);
      setLoaded(false);
      setError(null);
    }
  }, [userId, limit, load]);

  // Для графика — в хронологическом порядке (API отдаёт новые первыми).
  const chartData = [...rows].reverse().map((r) => ({
    t: fmt(r.timestamp),
    hp: r.hp,
    mp: r.mp,
    exp: r.exp,
    gp: r.gp,
  }));

  return (
    <section className="card">
      <h2>История статов</h2>

      {usersError && (
        <p data-testid="status-users-error" className="msg err">
          Не удалось загрузить список пользователей: {usersError}
        </p>
      )}

      <div className="row">
        <select
          data-testid="select-userid"
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
        >
          <option value="">— выберите пользователя —</option>
          {users.map((u) => (
            <option key={u.id} value={u.id}>
              {u.name} ({u.id})
            </option>
          ))}
        </select>
        <input
          data-testid="input-limit"
          type="number"
          min="0"
          value={limit}
          onChange={(e) => setLimit(e.target.value)}
          placeholder="лимит"
          style={{ maxWidth: 110 }}
        />
        <button
          data-testid="button-reload"
          type="button"
          onClick={() => load(userId, limit)}
          disabled={loading || !userId}
        >
          {loading ? "Загрузка…" : "Обновить"}
        </button>
      </div>

      {users.length === 0 && !usersError && (
        <p data-testid="no-users" className="hint">
          Пока нет пользователей. Добавьте пользователя в форме слева.
        </p>
      )}

      {error && (
        <p data-testid="status-error" className="msg err">
          {error}
        </p>
      )}

      {loaded && !error && rows.length === 0 && (
        <p data-testid="empty-state" className="hint">
          Нет данных. Снапшоты появятся после первого опроса (по умолчанию — раз в час).
        </p>
      )}

      {rows.length > 0 && (
        <>
          <div style={{ width: "100%", height: 300, marginTop: 16 }}>
            <ResponsiveContainer>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#2a2f45" />
                <XAxis dataKey="t" tick={{ fontSize: 11 }} minTickGap={24} />
                <YAxis tick={{ fontSize: 11 }} />
                <Tooltip
                  contentStyle={{
                    background: "#1b1f33",
                    border: "1px solid #2a2f45",
                  }}
                />
                <Legend />
                <Line type="monotone" dataKey="hp" stroke="#e74c3c" dot={false} />
                <Line type="monotone" dataKey="mp" stroke="#3498db" dot={false} />
                <Line type="monotone" dataKey="exp" stroke="#f1c40f" dot={false} />
                <Line type="monotone" dataKey="gp" stroke="#2ecc71" dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>

          <table data-testid="table-history" className="table">
            <thead>
              <tr>
                <th>Время</th>
                <th>HP</th>
                <th>MP</th>
                <th>Exp</th>
                <th>Gold</th>
                <th>Lvl</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r, i) => (
                <tr key={i} data-testid={`row-snapshot-${i}`}>
                  <td>{fmt(r.timestamp)}</td>
                  <td>{r.hp}</td>
                  <td>{r.mp}</td>
                  <td>{r.exp}</td>
                  <td>{r.gp}</td>
                  <td>{r.lvl}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </section>
  );
}
