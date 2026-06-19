import { useState } from "react";
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
import { getStatsHistory } from "./api";

const fmt = (ts) =>
  new Date(ts).toLocaleString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });

// Экран истории статов -> GET /api/v1/user/stats/history.
export default function StatsHistory() {
  const [userId, setUserId] = useState("");
  const [limit, setLimit] = useState("50");
  const [rows, setRows] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [loaded, setLoaded] = useState(false);

  async function load(e) {
    e && e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const data = await getStatsHistory(userId, Number(limit) || 0);
      setRows(data || []);
      setLoaded(true);
    } catch (err) {
      setError(err.message);
      setRows([]);
    } finally {
      setLoading(false);
    }
  }

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
      <form onSubmit={load} className="row">
        <input
          data-testid="input-userid"
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          placeholder="User ID"
          required
        />
        <input
          data-testid="input-limit"
          type="number"
          min="0"
          value={limit}
          onChange={(e) => setLimit(e.target.value)}
          placeholder="лимит"
          style={{ maxWidth: 110 }}
        />
        <button data-testid="button-load" type="submit" disabled={loading}>
          {loading ? "Загрузка…" : "Показать"}
        </button>
      </form>

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
