import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCreateAssistant } from '@/hooks/queries';
import { toApiError } from '@/api';
import { AssistantForm, type AssistantFormValues } from '@/components/AssistantForm';

export function NewAssistantPage() {
  const mutation = useCreateAssistant();
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const handleSubmit = async (values: AssistantFormValues) => {
    setServerError(null);
    try {
      const created = await mutation.mutateAsync({
        categoryId: values.categoryId,
        name: values.name,
        description: values.description,
        model: values.model,
        systemPrompt: values.systemPrompt,
        exampleUserPrompt: values.exampleUserPrompt || null,
        isActive: values.isActive,
      });
      // После создания открываем карточку — каталог уже инвалидирован хуком.
      navigate(`/assistants/${created.id}`);
    } catch (err) {
      setServerError(toApiError(err).message);
    }
  };

  return (
    <div className="form-narrow">
      <div className="title-row">
        <div className="title-row__text">
          <h1>Новый ассистент</h1>
          <p className="sub">Карточка каталога и правила запуска</p>
        </div>
      </div>
      <AssistantForm
        submitting={mutation.isPending}
        submitLabel="Создать ассистента"
        serverError={serverError}
        onSubmit={handleSubmit}
        onCancel={() => navigate(-1)}
      />
    </div>
  );
}
