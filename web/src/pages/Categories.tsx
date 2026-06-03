import { useEffect, useState, type ReactElement } from "react";
import { listCategories, deleteCategory, createCategory, type ProductCategory } from "../api/products";
import { useToast } from "../context/ToastContext";

export default function Categories() {
  const { showToast } = useToast();
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [newName, setNewName] = useState("");
  const [newDesc, setNewDesc] = useState("");
  const [newParent, setNewParent] = useState("");

  const load = () => {
    setLoading(true);
    listCategories().then(setCategories).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim()) return;
    try {
      await createCategory({ name: newName, description: newDesc || undefined, parent_id: newParent || undefined });
      setNewName("");
      setNewDesc("");
      setNewParent("");
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed", "error");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this category?")) return;
    try {
      await deleteCategory(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  const rootCats = categories.filter(c => !c.parent_id);
  const childMap = new Map<string, ProductCategory[]>();
  for (const c of categories) {
    if (c.parent_id) {
      const arr = childMap.get(c.parent_id) || [];
      arr.push(c);
      childMap.set(c.parent_id, arr);
    }
  }

  const renderCategory = (c: ProductCategory, depth = 0) => (
    <tr key={c.id}>
      <td style={{ paddingLeft: 20 * depth + 8 }}>{c.name}</td>
      <td>{c.description || "—"}</td>
      <td>{c.is_active ? "Yes" : "No"}</td>
      <td>
        <button onClick={() => handleDelete(c.id)}>Delete</button>
      </td>
    </tr>
  );

  const renderWithChildren = (cat: ProductCategory, depth = 0): ReactElement[] => {
    const rows = [renderCategory(cat, depth)];
    const children = childMap.get(cat.id) || [];
    for (const child of children) {
      rows.push(...renderWithChildren(child, depth + 1));
    }
    return rows;
  };

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Product Categories</h1>
      </div>

      <form onSubmit={handleCreate} style={{ display: "flex", gap: 8, maxWidth: "100%", marginBottom: 16, alignItems: "flex-end" }}>
        <label style={{ flex: 1 }}>Name
          <input value={newName} onChange={e => setNewName(e.target.value)} required />
        </label>
        <label style={{ flex: 1 }}>Description
          <input value={newDesc} onChange={e => setNewDesc(e.target.value)} />
        </label>
        <label>Parent
          <select value={newParent} onChange={e => setNewParent(e.target.value)}>
            <option value="">None (root)</option>
            {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        </label>
        <button type="submit">Add</button>
      </form>

      <table>
        <thead><tr><th>Name</th><th>Description</th><th>Active</th><th>Actions</th></tr></thead>
        <tbody>
          {rootCats.map(c => renderWithChildren(c))}
        </tbody>
      </table>
    </div>
  );
}
