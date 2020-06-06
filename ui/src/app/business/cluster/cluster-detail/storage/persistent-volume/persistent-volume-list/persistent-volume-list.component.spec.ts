import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeListComponent } from './persistent-volume-list.component';

describe('PersistentVolumeListComponent', () => {
  let component: PersistentVolumeListComponent;
  let fixture: ComponentFixture<PersistentVolumeListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
