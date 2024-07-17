import { indexOf } from "lodash";
import { FC, useCallback, useState } from "react";
import { v4 } from "uuid";
import { AnimatedButton } from "../../components/button";
import { CodeEditor } from "../../components/editor";
import { Icons } from "../../components/icons";
import { Loading } from "../../components/loading";
import { InternalPage } from "../../components/page";
import { Table } from "../../components/table";
import { InternalRoutes } from "../../config/routes";
import { DatabaseType, useRawExecuteLazyQuery } from "../../generated/graphql";
import { useAppSelector } from "../../store/hooks";
import classNames from "classnames";

type IRawExecuteCellProps = {
    cellId: string;
    onAdd: (cellId: string) => void;
    onDelete?: (cellId: string) => void;
    showTools?: boolean;
}

const RawExecuteCell: FC<IRawExecuteCellProps> = ({ cellId, onAdd, onDelete, showTools }) => {
    const [code, setCode] = useState("");
    const [rawExecute, { data: rows, loading, error }] = useRawExecuteLazyQuery();

    const current = useAppSelector(state => state.auth.current);

    const handleRawExecute = useCallback(() => {
        rawExecute({
            variables: {
                type: current?.Type as DatabaseType,
                query: code,
            },
        })
    }, [code, current?.Type, rawExecute]);

    const handleAdd = useCallback(() => {
        onAdd(cellId);
    }, [cellId, onAdd]);


    const handleDelete = useCallback(() => {
        onDelete?.(cellId);
    }, [cellId, onDelete]);

    return <div className="flex flex-col grow group/cell">
            <div className="relative">
                <div className="flex grow h-[150px] border border-gray-200 rounded-md overflow-hidden dark:bg-white/10 dark:border-white/5">
                    {
                        loading
                        ? <Loading hideText={true} />
                        : <CodeEditor language="sql" value={code} setValue={setCode} onRun={handleRawExecute} />
                    }
                </div>
                <div className={classNames("absolute -bottom-3 z-20 flex justify-between px-3 pr-8 w-full opacity-0 transition-all duration-500 group-hover/cell:opacity-100", {
                    "opacity-100": showTools,
                })}>
                    <div className="flex gap-2">
                        <AnimatedButton icon={Icons.PlusCircle} label="Add" onClick={handleAdd} />
                        {
                            onDelete != null &&
                            <AnimatedButton className="bg-red-100/80 hover:bg-red-200 dark:hover:bg-red-900" iconClassName="stroke-red-800" labelClassName="text-red-800"  icon={Icons.Delete} label="Delete" onClick={handleDelete} />
                        }
                    </div>
                    <AnimatedButton className="bg-green-200 hover:bg-green-400 dark:hover:bg-green-900" iconClassName="stroke-green-800" labelClassName="text-green-800" icon={Icons.CheckCircle} label="Submit query" onClick={handleRawExecute} />
                </div>
            </div>
            {
                error != null &&
                <div className="flex items-center justify-between mt-4">
                    <div className="text-sm text-red-500 w-[33vw]">{error?.message ?? ""}</div>
                </div>
            }
            {
                rows != null &&
                <div className="flex flex-col w-full h-[250px] mt-4">
                    <Table columns={rows.RawExecute.Columns.map(c => c.Name)} columnTags={rows.RawExecute.Columns.map(c => c.Type)}
                        rows={rows.RawExecute.Rows} totalPages={1} currentPage={1} disableEdit={true} />
                </div>
            }
        </div>
}

export const RawExecutePage: FC = () => {
    const [cellIds, setCellIds] = useState<string[]>([v4()]);
    
    const handleAdd = useCallback((id: string) => {
        const index = indexOf(cellIds, id);
        const newCellIds = [...cellIds];
        newCellIds.splice(index+1, 0, v4());
        setCellIds(newCellIds);
    }, [cellIds]);

    const handleDelete = useCallback((cellId: string) => {
        if (cellIds.length <= 1) {
            return;
        }
        setCellIds(ids => ids.filter(id => id !== cellId));
    }, [cellIds.length]);

    return (
        <InternalPage routes={[InternalRoutes.RawExecute]}>
            <div className="flex justify-center items-center w-full">
                <div className="w-full max-w-[1000px] flex flex-col gap-4">
                    {
                        cellIds.map((cellId, index) => (
                            <>
                                {index > 0 && <div className="border-dashed border-t border-gray-300 my-2 dark:border-neutral-600"></div>}
                                <RawExecuteCell key={cellId} cellId={cellId} onAdd={handleAdd} onDelete={cellIds.length <= 1 ? undefined : handleDelete}
                                    showTools={cellIds.length === 1} />
                            </>
                        ))
                    }
                </div>
            </div>
        </InternalPage>
    )
}   