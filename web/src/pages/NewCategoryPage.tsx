import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCreateCategory } from '@/hooks/queries';
import { toApiError } from '@/api';
import { Alert } from '@/components/States';
import { Button, FieldError } from '@/components/Button';
import { required } from '@/lib/validation';

export function NewCategoryPage() {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [errors, setErrors] = useState<{ name?: string }>({});
  const [serverError, setServerError] = useState<string | null>(null);

  const mutation = useCreateCategory();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setServerError(null);
    const nameError = required(name, 'Название');
    setErrors({ name: nameError });
    if (nameError) return;

    try {
      await mutation.mutateAsync({ name: name.trim(), description: description.trim() || null });
      navigate('/assistants');
    } catch (err) {
      setServerError(toApiError(err).message);
    }
  };

  return (
    <div className="form-narrow">
      <div className="title-row">
        <div className="title-row__text">
          <h1>Новая категория</h1>
          <p className="sub">Предметная область для ассистентов</p>
        </div>
      </div>

      <form className="form" onSubmit={handleSubmit} noValidate>
        {serverError && <Alert variant="error">{serverError}</Alert>}

        <div className={`field ${errors.name ? 'field--invalid' : ''}`}>
          <label className="field__label" htmlFor="name">
            Название *
          </label>
          <input
            id="name"
            className="input"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Например, Еда"
          />
          <FieldError>{errors.name}</FieldError>
        </div>

        <div className="field">
          <label className="field__label" htmlFor="description">
            Описание
          </label>
          <textarea
            id="description"
            className="input"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Необязательное описание категории"
          />
        </div>

        <div className="form-actions">
          <Button type="submit" variant="primary" disabled={mutation.isPending}>
            {mutation.isPending ? 'Создаём…' : 'Создать категорию'}
          </Button>
          <Button type="button" variant="ghost" onClick={() => navigate(-1)}>
            Отмена
          </Button>
        </div>
      </form>
    </div>
  );
}
