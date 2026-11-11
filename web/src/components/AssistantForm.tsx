import { useEffect, useState } from 'react';
import { useCategories, useModels } from '@/hooks/queries';
import { Button, FieldError } from '@/components/Button';
import { isUuid, required } from '@/lib/validation';
import type { Assistant } from '@/types/api';

export interface AssistantFormValues {
  categoryId: string;
  name: string;
  description: string;
  model: string;
  systemPrompt: string;
  exampleUserPrompt: string;
  isActive: boolean;
}

type Errors = Partial<Record<keyof AssistantFormValues, string>>;

interface AssistantFormProps {
  initial?: Assistant;
  submitting: boolean;
  submitLabel: string;
  serverError?: string | null;
  onSubmit: (values: AssistantFormValues) => void | Promise<void>;
  onCancel: () => void;
}

const EMPTY: AssistantFormValues = {
  categoryId: '',
  name: '',
  description: '',
  model: '',
  systemPrompt: '',
  exampleUserPrompt: '',
  isActive: true,
};

export function AssistantForm({
  initial,
  submitting,
  submitLabel,
  serverError,
  onSubmit,
  onCancel,
}: AssistantFormProps) {
  const categoriesQuery = useCategories();
  const modelsQuery = useModels();
  const [values, setValues] = useState<AssistantFormValues>(() =>
    initial
      ? {
          categoryId: initial.categoryId,
          name: initial.name,
          description: initial.description,
          model: initial.model,
          systemPrompt: initial.systemPrompt ?? '',
          exampleUserPrompt: initial.exampleUserPrompt ?? '',
          isActive: initial.isActive,
        }
      : EMPTY,
  );
  const [errors, setErrors] = useState<Errors>({});
  const [modelSearch, setModelSearch] = useState('');

  const models = modelsQuery.data ?? [];
  const hasModels = modelsQuery.isSuccess && models.length > 0;
  // Текущая модель может отсутствовать в списке провайдера (например, старая запись).
  const modelOptions =
    values.model && !models.includes(values.model) ? [values.model, ...models] : models;

  const normalizedModelSearch = modelSearch.trim().toLowerCase();
  const filteredModels = modelOptions.filter((model) =>
    model.toLowerCase().includes(normalizedModelSearch),
  );
  // Текущее значение всегда должно оставаться в options, иначе select покажет пустоту.
  const visibleModels =
    values.model && !filteredModels.includes(values.model)
      ? [values.model, ...filteredModels]
      : filteredModels;

  // При создании подставляем первую доступную модель, когда список загрузился.
  useEffect(() => {
    const available = modelsQuery.data;
    if (initial || !available?.length) return;
    if (!available.includes(values.model)) {
      setValues((prev) => ({ ...prev, model: available[0] }));
    }
  }, [initial, modelsQuery.data, values.model]);

  const set = <K extends keyof AssistantFormValues>(key: K, value: AssistantFormValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  const validate = (): Errors => {
    const next: Errors = {};
    next.name = required(values.name, 'Название');
    next.description = required(values.description, 'Описание');
    next.model = required(values.model, 'Модель');
    next.systemPrompt = required(values.systemPrompt, 'Системный промпт');
    if (!values.categoryId) next.categoryId = 'Выберите категорию';
    else if (!isUuid(values.categoryId)) next.categoryId = 'Некорректный идентификатор категории';
    return next;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const next = validate();
    setErrors(next);
    if (Object.values(next).some(Boolean)) return;
    onSubmit({
      ...values,
      name: values.name.trim(),
      description: values.description.trim(),
      model: values.model.trim(),
      systemPrompt: values.systemPrompt.trim(),
      exampleUserPrompt: values.exampleUserPrompt.trim(),
    });
  };

  const categories = categoriesQuery.data ?? [];

  return (
    <form className="form" onSubmit={handleSubmit} noValidate>
      <div className={`field ${errors.name ? 'field--invalid' : ''}`}>
        <label className="field__label" htmlFor="name">
          Название *
        </label>
        <input
          id="name"
          className="input"
          value={values.name}
          onChange={(e) => set('name', e.target.value)}
        />
        <FieldError>{errors.name}</FieldError>
      </div>

      <div className={`field ${errors.description ? 'field--invalid' : ''}`}>
        <label className="field__label" htmlFor="description">
          Описание *
        </label>
        <textarea
          id="description"
          className="input"
          value={values.description}
          onChange={(e) => set('description', e.target.value)}
        />
        <FieldError>{errors.description}</FieldError>
      </div>

      <div className={`field ${errors.categoryId ? 'field--invalid' : ''}`}>
        <label className="field__label" htmlFor="categoryId">
          Категория *
        </label>
        <select
          id="categoryId"
          className="input"
          value={values.categoryId}
          onChange={(e) => set('categoryId', e.target.value)}
          disabled={categoriesQuery.isLoading}
        >
          <option value="">— выберите категорию —</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
        <FieldError>{errors.categoryId}</FieldError>
        {categories.length === 0 && !categoriesQuery.isLoading && (
          <span className="field__hint">
            Нет доступных категорий. Сначала создайте категорию.
          </span>
        )}
      </div>

      <div className={`field ${errors.model ? 'field--invalid' : ''}`}>
        <label className="field__label" htmlFor="model">
          Модель *
        </label>
        {hasModels ? (
          <>
            <input
              className="input"
              aria-label="Поиск модели"
              placeholder="Поиск модели…"
              value={modelSearch}
              onChange={(e) => setModelSearch(e.target.value)}
            />
            <select
              id="model"
              className="input mono"
              value={values.model}
              onChange={(e) => set('model', e.target.value)}
            >
              {visibleModels.map((model) => (
                <option key={model} value={model}>
                  {model}
                </option>
              ))}
            </select>
            {filteredModels.length === 0 && (
              <span className="field__hint">
                По запросу ничего не найдено — показана текущая модель.
              </span>
            )}
          </>
        ) : (
          <input
            id="model"
            className="input mono"
            value={values.model}
            onChange={(e) => set('model', e.target.value)}
            placeholder={modelsQuery.isLoading ? 'Загружаем модели…' : 'mock-smart'}
            disabled={modelsQuery.isLoading}
          />
        )}
        <FieldError>{errors.model}</FieldError>
        {modelsQuery.isError && (
          <span className="field__hint">
            Не удалось загрузить список моделей — введите название вручную.
          </span>
        )}
        {modelsQuery.isSuccess && models.length === 0 && (
          <span className="field__hint">
            Провайдер не вернул модели — введите название вручную.
          </span>
        )}
        {hasModels && initial && !models.includes(values.model) && (
          <span className="field__hint">Текущая модель отсутствует в списке провайдера.</span>
        )}
      </div>

      <div className={`field ${errors.systemPrompt ? 'field--invalid' : ''}`}>
        <label className="field__label" htmlFor="systemPrompt">
          Системный промпт *
        </label>
        <textarea
          id="systemPrompt"
          className="input"
          rows={5}
          value={values.systemPrompt}
          onChange={(e) => set('systemPrompt', e.target.value)}
          placeholder="Правила поведения ассистента…"
        />
        <FieldError>{errors.systemPrompt}</FieldError>
      </div>

      <div className="field">
        <label className="field__label" htmlFor="exampleUserPrompt">
          Пример пользовательского промпта
        </label>
        <textarea
          id="exampleUserPrompt"
          className="input"
          value={values.exampleUserPrompt}
          onChange={(e) => set('exampleUserPrompt', e.target.value)}
          placeholder="курица, рис, томаты, сыр"
        />
        <span className="field__hint">
          Подсказка, которую увидит пользователь на карточке ассистента.
        </span>
      </div>

      <label className="check">
        <input
          type="checkbox"
          checked={values.isActive}
          onChange={(e) => set('isActive', e.target.checked)}
        />
        Ассистент активен
      </label>

      {serverError && <div className="alert alert--error">{serverError}</div>}

      <div className="form-actions">
        <Button type="submit" variant="primary" disabled={submitting}>
          {submitting ? 'Сохраняем…' : submitLabel}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel}>
          Отмена
        </Button>
      </div>
    </form>
  );
}
