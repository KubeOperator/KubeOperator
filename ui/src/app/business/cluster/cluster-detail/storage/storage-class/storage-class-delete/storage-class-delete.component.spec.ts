import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageClassDeleteComponent } from './storage-class-delete.component';

describe('StorageClassDeleteComponent', () => {
  let component: StorageClassDeleteComponent;
  let fixture: ComponentFixture<StorageClassDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ StorageClassDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageClassDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
