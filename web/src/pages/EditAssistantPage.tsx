import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAssistant, useUpdateAssistant } from '@/hooks/queries';
import { toApiError } from '@/api';
import { ErrorState, LoadingState } from '@/components/States';
import { AssistantForm, type AssistantFormValues } from '@/components/AssistantForm';

export function EditAssistantPage() {
  const { id } = useParams<{ id: string }>();
  const query = useAssistant(id);
  const mutation = useUpdateAssistant();
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  if (query.isLoading) return <LoadingState label="Загружаем ассистента…" />;
  if (query.isError) {
    return <ErrorState message={toApiError(query.error).message} onRetry={() => query.refetch()} />;
  }
  if (!query.data || !id) return null;

  const handleSubmit = async (values: AssistantFormValues) => {
    setServerError(null);
    try {
      await mutation.mutateAsync({
        id,
        payload: {
          categoryId: values.categoryId,
          name: values.name,
          description: values.description,
          model: values.model,
          systemPrompt: values.systemPrompt,
          exampleUserPrompt: values.exampleUserPrompt || null,
          isActive: values.isActive,
        },
      });
      navigate(`/assistants/${id}`);
    } catch (err) {
      setServerError(toApiError(err).message);
    }
  };

  return (
    <div className="form-narrow">
      <div className="title-row">
        <div className="title-row__text">
          <h1>{query.data.name}</h1>
          <p className="sub">Редактирование ассистента</p>
        </div>
      </div>
      <AssistantForm
        initial={query.data}
        submitting={mutation.isPending}
        submitLabel="Сохранить изменения"
        serverError={serverError}
        onSubmit={handleSubmit}
        onCancel={() => navigate(-1)}
      />
    </div>
  );
}
