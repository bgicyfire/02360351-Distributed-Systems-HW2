import {Injectable} from '@angular/core';
import {map, Observable, tap} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {Scooter} from '../dtos/scooter';
import {environment} from '../../environments/environment';
import {ToastService} from './toast.service';



@Injectable({
  providedIn: 'root'
})
export class ScootersService {
  private BASE_URL = environment.scootersApiUrl;

  constructor(private http: HttpClient, private toastService: ToastService) {
  }

  public getServersList(): Observable<{ servers: string[], responder: string, myLeader: string }> {
    return this.http.get<{ servers: string[], responder: string, myLeader: string }>(this.BASE_URL + '/servers');
  }

  public getScootersList(): Observable<Scooter[]> {
    return this.http.get<{scootersList:Scooter[], responder: string, myLeader: string}>(this.BASE_URL + '/scooters').pipe(
      tap(res=>this.toastService.showSuccess('Responder: ' + res.responder + ' his leader is: ' + res.myLeader)),
      map(res => res.scootersList.sort((a, b) => a.id.localeCompare(b.id)))
    );
  }

  public createScooter(id: string): Observable<Scooter> {
    let data: Scooter = {
      'id': id,
      'is_available': true,
      'total_distance': 0,
      'current_reservation_id': ''
    };
    return this.http.put<{newScooter:Scooter, responder: string, myLeader: string}>(this.BASE_URL + '/scooters/' + id, data).pipe(
      tap(res=>this.toastService.showSuccess('Responder: ' + res.responder + ' his leader is: ' + res.myLeader)),
      map(res => res.newScooter)
    );
  }

  public reserveScooter(scooterId: string): Observable<string> {
    return this.http.post<string>(this.BASE_URL + '/scooters/' + scooterId + '/reservations', {});
  }

  public releaseScooter(scooterId: string, reservationId: string, rideDistance: number): Observable<any> {
    // TODO: decide if also want to require the reservation id
    let data = {
      scooterId: scooterId,
      reservation_id: reservationId,
      ride_distance: rideDistance
    }
    return this.http.post<any>(this.BASE_URL + '/scooters/' + scooterId + '/releases', data);
  }
}
