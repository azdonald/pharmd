import { useEffect, useState, type ReactElement } from "react";
import { listCategories, deleteCategory, createCategory, type ProductCategory } from "../api/products";
import { useToast } from "../context/ToastContext";
import { Icon, PageHeader, Panel } from "../components/AdminComponents";


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
      setNewName(""); setNewDesc(""); setNewParent("");
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
    <tr key={c.id} className="hover:bg-surface-container-high/20 transition-colors">
      <td className="px-6 py-4 font-semibold text-on-surface" style={{ paddingLeft: 20 * depth + 24 }}>
        {depth > 0 && <Icon name="subdirectory_arrow_right" className="text-on-surface-variant mr-1 text-sm" />}
        {c.name}
      </td>
      <td className="px-6 py-4 text-body-md text-on-surface-variant">{c.description || "—"}</td>
      <td className="px-6 py-4 text-center">
        <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
          c.is_active ? "bg-secondary-container/20 text-secondary" : "bg-outline-variant/20 text-on-surface-variant"
        }`}>{c.is_active ? "Active" : "Inactive"}</span>
      </td>
      <td className="px-6 py-4 text-right">
        <button onClick={() => handleDelete(c.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
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

  return (
    <div>
      <PageHeader
        title="Product Categories"
        description="Organize products into categories and subcategories"
      />

      {/* Add form */}
      <form onSubmit={handleCreate} className="mb-8 p-4 rounded-xl border border-outline-variant bg-surface-container-lowest">
        <h3 className="font-semibold text-on-surface mb-3">Add New Category</h3>
        <div className="flex gap-3 items-end">
          <div className="flex-1">
            <label className="block text-xs font-medium text-on-surface-variant mb-1">Name *</label>
            <input value={newName} onChange={e => setNewName(e.target.value)} required
              className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
          </div>
          <div className="flex-1">
            <label className="block text-xs font-medium text-on-surface-variant mb-1">Description</label>
            <input value={newDesc} onChange={e => setNewDesc(e.target.value)}
              className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
          </div>
          <div>
            <label className="block text-xs font-medium text-on-surface-variant mb-1">Parent</label>
            <select value={newParent} onChange={e => setNewParent(e.target.value)}
              className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
              <option value="">None (root)</option>
              {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </div>
          <button type="submit"
            className="px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container transition-all">Add</button>
        </div>
      </form>

      {/* Table */}
      <Panel>
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading categories...</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Description</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {rootCats.map(c => renderWithChildren(c))}
              </tbody>
            </table>
          )}
        </div>
      </Panel>
    </div>
  );
}
