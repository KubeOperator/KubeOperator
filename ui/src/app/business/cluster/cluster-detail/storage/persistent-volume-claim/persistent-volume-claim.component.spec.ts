import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeClaimComponent } from './persistent-volume-claim.component';

describe('PersistentVolumeClaimComponent', () => {
  let component: PersistentVolumeClaimComponent;
  let fixture: ComponentFixture<PersistentVolumeClaimComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeClaimComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeClaimComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
