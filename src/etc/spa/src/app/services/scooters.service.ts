import {Injectable} from '@angular/core';
import {map, Observable, of} from 'rxjs';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ScootersService {
  constructor(private http: HttpClient) {
  }

  public getServersList(): Observable<string[]> {
    return this.http.get<{servers:string[]}>('http://localhost:50053/server')
      .pipe(map(a=>a.servers));
  }
}
