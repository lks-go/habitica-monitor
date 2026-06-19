import { useState } from "react";
import { createUser } from "./api";

// Форма добавления пользователя -> POST /api/v1/user.
export default function AddUser({ onCreated }) {
  const [form, setForm] = useState({ id: "", api_key: "", name: "" });
  const [status, setStatus] = useState(null); // {type, msg}
  const [busy, setBusy] = useState(false);

  const update = (k) => (e) => setForm({ ...form, [k]: e.target.value });

  async function submit(e) {
    e.preventDefault();
    setBusy(true);
    setStatus(null);
    try {
      const u = await createUser(form);
      setStatus({ type: "ok", msg: `Пользователь «${u.name}» добавлен` });
      setForm({ id: "", api_key: "", name: "" });
      onCreated && onCreated(u);
    } catch (err) {
      setStatus({ type: "err", msg: err.message });
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h2>Добавить пользователя</h2>
      <form onSubmit={submit} className="form">
        <label>
          <span>Имя</span>
          <input
            data-testid="input-name"
            value={form.name}
            onChange={update("name")}
            placeholder="alice"
            required
          />
        </label>
        <label>
          <span>User ID (x-api-user)</span>
          <input
            data-testid="input-id"
            value={form.id}
            onChange={update("id")}
            placeholder="abcd-1234-..."
            required
          />
        </label>
        <label>
          <span>API Token (x-api-key)</span>
          <input
            data-testid="input-apikey"
            type="password"
            value={form.api_key}
            onChange={update("api_key")}
            placeholder="секретный токен"
            required
          />
        </label>
        <button data-testid="button-submit" type="submit" disabled={busy}>
          {busy ? "Сохранение…" : "Добавить"}
        </button>
      </form>
      {status && (
        <p
          data-testid="status-message"
          className={status.type === "ok" ? "msg ok" : "msg err"}
        >
          {status.msg}
        </p>
      )}
      <p className="hint">
        User ID и API Token берутся в Habitica: Settings → API.
      </p>
    </section>
  );
}
