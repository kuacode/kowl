import { observer } from "mobx-react";
import { PageComponent, PageInitHelper } from "../Page";
import { Table, Tag } from 'antd';
import { Key, useEffect, useState } from "react";
import { api } from '../../../state/backendApi';
import { appGlobal } from '../../../state/appGlobal';

const MyTable = ({
    selected = '',
    ds = [],
    changeSelected = (id: string) => {
        // no-op default
    }
}: {
    selected?: string,
    ds: any[],
    changeSelected?: (id: string) => void
}) => {
    // 列定义
    const columns = [
        { title: 'Name', dataIndex: 'name', key: 'name' },
        {
            title: 'Brokers', dataIndex: 'brokers', key: 'brokers',
            render: (brokers: string[]) => (
                <>
                    {brokers.map(broker => (
                        <Tag color="blue" key={broker}>
                            {broker}
                        </Tag>
                    ))}
                </>
            )
        }
    ];

    const [selectedRowKeys, setSelectedRowKeys] = useState<Key[]>([]);

    const rowSelection = {
        type: 'radio' as const,
        columnWidth: 48,
        selectedRowKeys,
        onChange: (keys: Key[]) => {
            setSelectedRowKeys(keys);
            changeSelected(keys[0] as string);
        }
    }

    useEffect(() => {
        if (selected) {
            setSelectedRowKeys([selected]);
        }
    }, [selected]);

    return (
        <Table
            rowKey='id'
            columns={columns}
            dataSource={ds}
            rowSelection={rowSelection}
        />
    );
};

@observer
class ClusterList extends PageComponent {

    initPage(p: PageInitHelper): void {
        p.title = 'Home';
        p.addBreadcrumb('Home', '/clusters');

        this.refreshData(false);
        appGlobal.onRefresh = () => this.refreshData(true);
    }

    render() {
        const clusterData: any = api.clusterData || [];
        return <>
            <MyTable ds={clusterData.data || []} selected={clusterData.selected} changeSelected={this.setCurrentCluster} />
        </>;
    }

    refreshData(force: boolean) {
        api.refreshClusterList(force);
    }

    setCurrentCluster(id: string) {
        api.setCurrentCluster(id);
    }
}

export default ClusterList;
