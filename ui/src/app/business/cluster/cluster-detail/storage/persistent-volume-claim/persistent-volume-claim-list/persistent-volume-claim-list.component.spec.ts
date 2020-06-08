import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeClaimListComponent } from './persistent-volume-claim-list.component';

describe('PersistentVolumeClaimListComponent', () => {
  let component: PersistentVolumeClaimListComponent;
  let fixture: ComponentFixture<PersistentVolumeClaimListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeClaimListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeClaimListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
