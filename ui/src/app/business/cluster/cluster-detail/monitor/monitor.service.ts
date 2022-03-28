import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class MonitorService {
    baseUrl = '';
    
    cpuCommon = 'api/v1/query_range?query=sum by(instance)(irate(node_cpu_seconds_total{instance="{instance}", mode={mode}}[5m])) * 100&start={start}&end={end}&step=300'
    
    memeryTotal = 'api/v1/query_range?query=node_memory_MemTotal_bytes{instance="{instance}"}&start={start}&end={end}&step=300'
    memeryUsed = 'api/v1/query_range?query=node_memory_MemTotal_bytes{instance="{instance}"} - node_memory_MemFree_bytes{instance="{instance}"} - (node_memory_Cached_bytes{instance="{instance}"}%2Bnode_memory_Buffers_bytes{instance="{instance}"})&start={start}&end={end}&step=300'
    memeryCacheBuffer = 'api/v1/query_range?query=node_memory_Cached_bytes{instance="{instance}"}%2Bnode_memory_Buffers_bytes{instance="{instance}"}&start={start}&end={end}&step=300'
    memeryFree = 'api/v1/query_range?query=node_memory_MemFree_bytes{instance="{instance}"}&start={start}&end={end}&step=300'
    memerySWAPUsed = 'api/v1/query_range?query=(node_memory_SwapTotal_bytes{instance="{instance}"} - node_memory_SwapFree_bytes{instance="{instance}"})&start={start}&end={end}&step=300'
    
    disk = 'api/v1/query_range?query=100 - ((node_filesystem_avail_bytes{instance="{instance}",device!~"rootfs"} * 100) / node_filesystem_size_bytes{instance="{instance}",device!~"rootfs"})&start={start}&end={end}&step=300'
    
    networkRecv = 'api/v1/query_range?query=irate(node_network_receive_bytes_total{instance="{instance}"}[5m])*8&start={start}&end={end}&step=300'
    networkTrans = 'api/v1/query_range?query=irate(node_network_transmit_bytes_total{instance="{instance}"}[5m])*8&start={start}&end={end}&step=300'
    
    constructor(private http: HttpClient) {}

    QueryCPU(node: string,mode: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.cpuCommon.replace(/{instance}/g, node).replace(/{mode}/g, mode).replace('{start}', start).replace('{end}', end));
    }

    QueryMemeryTotal(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.memeryTotal.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }
    QueryMemeryUsed(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.memeryUsed.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }
    QueryMemeryCacheBuffer(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.memeryCacheBuffer.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }
    QueryMemeryFree(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.memeryFree.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }
    QueryMemerySWAPUsed(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.memerySWAPUsed.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }

    QueryDisk(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.disk.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }

    QueryNetworkRecv(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.networkRecv.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }
    QueryNetworkTrans(node: string, start: string, end: string): Observable<any> {
        return this.http.get<any>(this.baseUrl + this.networkTrans.replace(/{instance}/g, node).replace('{start}', start).replace('{end}', end));
    }
}