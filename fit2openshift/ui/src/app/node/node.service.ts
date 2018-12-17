import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {Node} from './node';

const baseNodeUrl = 'api/v1/cluster/{clusterId}/node/';

@Injectable({
  providedIn: 'root'
})
export class NodeService {

  constructor(private http: HttpClient) {
  }

  listNodes(clusterId): Observable<Node[]> {
    return this.http.get<Node[]>(baseNodeUrl.replace('{clusterId}', clusterId)).pipe(
      catchError(err => throwError(err))
    );
  }

  createNode(clusterId, node: Node): Observable<Node> {
    return this.http.post<Node>(baseNodeUrl.replace('{clusterId}', clusterId), node).pipe(
      catchError(err => throwError(err))
    );
  }

  deleteNode(clusterId, nodeId): Observable<any> {
    return this.http.delete(`${baseNodeUrl.replace('{clusterId}', clusterId)}/${nodeId}`).pipe(
      catchError(err => throwError(err))
    );
  }
}
